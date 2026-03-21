package message

import (
	"context"
	"mx-shop-api/user-interation-web/forms"
	"mx-shop-api/user-interation-web/global"
	"mx-shop-api/user-interation-web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Create(ctx *gin.Context) {
	s := zap.S()
	uidStr, _ := ctx.Get("userId")
	uid := int32(uidStr.(uint))
	var messageForm forms.MessageReqForm
	err := ctx.ShouldBind(&messageForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := global.MessageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      uid,
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
		"id":  resp.Id,
	})
}

func List(ctx *gin.Context) {
	s := zap.S()
	uidStr, _ := ctx.Get("userId")
	uid := int32(uidStr.(uint))
	resp, err := global.MessageClient.MessageList(context.Background(), &proto.MessageRequest{
		UserId: uid,
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := make(map[string]interface{})
	data := make([]interface{}, 0)
	response["total"] = resp.Total
	/**
	{
		total:0,
		data:[
			{
				messageType:0,
				subject:"",
				message:"",
				file:""
			}
		]
	}
	*/
	for _, msg := range resp.Data {
		data = append(data, &proto.MessageResponse{
			Id:          msg.Id,
			UserId:      msg.UserId,
			MessageType: msg.MessageType,
			Subject:     msg.Subject,
			Message:     msg.Message,
			File:        msg.File,
		})
	}
	response["data"] = data
	ctx.JSON(http.StatusOK, response)
}
