package api

import (
	"context"
	"fmt"
	"math/rand"
	"mx-shop-api/user-web/forms"
	"mx-shop-api/user-web/global"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

// 生成 width 宽度的验证码
func GenSmsCode(width int) string {
	numeric := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSMS(ctx *gin.Context) {
	s := zap.S()
	sendSmsForm := forms.SendSmsForm{}
	err := ctx.ShouldBind(&sendSmsForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecret)
	// if err != nil {
	// 	panic(err)
	// }
	mobile := sendSmsForm.Mobile
	smsCode := GenSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = mobile                        //手机号
	request.QueryParams["SignName"] = ""                                //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = ""                            //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	// response, err := client.ProcessCommonRequest(request)
	// fmt.Print(client.DoAction(request, response))
	// if err != nil {
	// 	fmt.Print(err.Error())
	// }
	s.Infof("mobile:%s,type:%d,sms_code:%s", sendSmsForm.Mobile, sendSmsForm.Type, smsCode)
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

	Rdb.Set(Ctx, mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second) // 5分钟过期

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
