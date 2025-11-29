package api

import (
	"context"
	"fmt"
	"mx-shop-api/user-web/forms"
	"mx-shop-api/user-web/global"
	"mx-shop-api/user-web/global/response"
	"mx-shop-api/user-web/middlewares"
	"mx-shop-api/user-web/models"
	"mx-shop-api/user-web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			switch status.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": status.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "internal error",
				})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "params error",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	s := zap.S()

	// 拨号连接user grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConf.Host, global.ServerConfig.UserSrvConf.Port), grpc.WithInsecure())
	if err != nil {
		s.Errorw("connect to user service error:", err.Error())
		return
	}

	// 生成grpc的client并调用接口
	userClient := proto.NewUserClient(userConn)
	pn := ctx.DefaultQuery("pn", "1")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	pageInfo := proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	}
	rsp, err := userClient.GetUserList(context.Background(), &pageInfo)
	if err != nil {
		s.Errorw("GetUserList Err:", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		data := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Mobile:   value.Mobile,
			Gender:   value.Gender,
			BirthDay: response.JsonTime(time.Unix(int64(value.BirthDay), 0)), // 将 uint64 类型转为 time.Time 类型
		}
		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)
}

func PasswordLogin(ctx *gin.Context) {
	s := zap.S()
	loginForm := forms.PasswordLoginForm{}

	err := ctx.ShouldBind(&loginForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// 对用户传入验证码进行验证
	if verifyResult := store.Verify(loginForm.CaptchaId, loginForm.Captcha, true); !verifyResult {
		s.Infof("CaptchaId:%s , Captcha:%s", loginForm.CaptchaId, loginForm.Captcha)
		// 验证不通过
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConf.Host, global.ServerConfig.UserSrvConf.Port), grpc.WithInsecure())
	if err != nil {
		s.Errorw("connect to user service error:", err.Error())
		return
	}

	userClient := proto.NewUserClient(userConn)
	// 登录逻辑验证
	// 首先查询是否有该用户
	resp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: loginForm.Mobile,
	})
	// 没查到用户
	if err != nil {
		statusInfo, ok := status.FromError(err)
		if ok {
			switch statusInfo.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录错误",
				})

			}
		}
		return
	}
	// 查到用户了，应当检查用户传入的密码与存储的密码是否一致
	checkResp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          loginForm.Password, // 入参传的是明文
		EncryptedPassword: resp.Password,      // 库里存的是密文
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "登录失败",
		})
		return
	}
	if checkResp.Success {
		// 生成token
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(resp.Id),
			NickName:    resp.NickName,
			AuthorityId: uint(resp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),               // 签名的生效时间
				ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期
				Issuer:    "lele",                          // 哪个机构的认证签名
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "内部生成token错误",
			})
			return
		}

		// 密码正确
		ctx.JSON(http.StatusOK, gin.H{
			"msg":       "登录成功",
			"id":        resp.Id,
			"nickname":  resp.NickName,
			"token":     token,
			"expire_at": (time.Now().Unix() + 60*60*24*30) * 1000,
		})
	} else {
		// 密码错误
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "密码错误",
		})
	}
}

func Register(ctx *gin.Context) {
	s := zap.S()
	registerForm := forms.RegisterForm{}
	err := ctx.ShouldBind(&registerForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvConf.Host, global.ServerConfig.UserSrvConf.Port), grpc.WithInsecure())
	if err != nil {
		s.Errorw("connect to user service error:", err.Error())
		return
	}

	userClient := proto.NewUserClient(userConn)
	// 验证验证码是否正确
	// 将验证码保存到redis
	Rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		Password: "",
		DB:       0, // 默认DB
	})

	_, err = Rdb.Ping(Ctx).Result()
	if err != nil {
		s.Errorf("Redis 连接失败: %v", err)
	}
	s.Infoln("Redis 连接成功!")

	val, _ := Rdb.Get(Ctx, registerForm.Mobile).Result()
	// if err == redis.Nil {
	// 	s.Error("Redis value not found")
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"msg": "验证码错误",
	// 	})
	// 	return
	// } else if val != registerForm.Code {
	// 	s.Error("Redis value not found")
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"msg": "验证码错误",
	// 	})
	// 	return
	// }

	s.Infof("Rdb.Get>>> %s", val)

	// 创建用户 + 登录成功（返回token）
	resp, err := userClient.CreateUser(Ctx, &proto.CreateUserInfo{
		Mobile:   registerForm.Mobile,
		Password: registerForm.Password,
		NickName: "新用户" + registerForm.Mobile,
	})
	if err != nil {
		s.Error("创建用户失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "创建用户失败",
			"err": err.Error(),
		})
		return
	}

	// 生成token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(resp.Id),
		NickName:    resp.NickName,
		AuthorityId: uint(resp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期
			Issuer:    "lele",                          // 哪个机构的认证签名
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "内部生成token错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":       "注册并登录成功",
		"id":        resp.Id,
		"nickname":  resp.NickName,
		"token":     token,
		"expire_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
