package api

import (
	"context"
	"fmt"
	"mx-shop-api/user-web/forms"
	"mx-shop-api/user-web/global"
	"mx-shop-api/user-web/global/response"
	"mx-shop-api/user-web/proto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
