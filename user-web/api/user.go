package api

import (
	"context"
	"fmt"
	"mx-shop-api/user-web/proto"
	"net/http"

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
	ip := "127.0.0.1"
	port := 50051

	// 拨号连接user grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		s.Errorw("connect to user service error:", err.Error())
		return
	}

	// 生成grpc的client并调用接口
	userClient := proto.NewUserClient(userConn)
	pageInfo := proto.PageInfo{
		Pn:    0,
		PSize: 0,
	}
	rsp, err := userClient.GetUserList(context.Background(), &pageInfo)
	if err != nil {
		s.Errorw("GetUserList Err:", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		data := make(map[string]interface{})
		data["id"] = value.Id
		data["name"] = value.NickName
		data["mobile"] = value.Mobile
		data["gender"] = value.Gender
		data["birthday"] = value.BirthDay

		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)
}
