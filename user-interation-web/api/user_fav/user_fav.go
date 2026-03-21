package userfav

import (
	"context"
	"mx-shop-api/user-interation-web/forms"
	"mx-shop-api/user-interation-web/global"
	"mx-shop-api/user-interation-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Create(ctx *gin.Context) {
	s := zap.S()
	uidStr, _ := ctx.Get("userId")
	uid := int32(uidStr.(uint))
	var userFavForm forms.UserFavReqForm
	err := ctx.ShouldBind(&userFavForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = global.UserFavClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  uid,
		GoodsId: userFavForm.GoodsId,
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func Del(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	uidStr, _ := ctx.Get("userId")
	uid := int32(uidStr.(uint))

	_, err := global.UserFavClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  uid,
		GoodsId: int32(id),
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func List(ctx *gin.Context) {
	s := zap.S()
	uidStr, _ := ctx.Get("userId")
	uid := int32(uidStr.(uint))
	resp, err := global.UserFavClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: uid,
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	data := make([]interface{}, 0)
	for _, v := range resp.Data {
		data = append(data, &proto.UserFavResponse{
			UserId:  v.UserId,
			GoodsId: v.GoodsId,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  data,
	})
}

func Get(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	uidStr, _ := ctx.Get("userId")
	uid := int32(uidStr.(uint))
	_, err := global.UserFavClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  uid,
		GoodsId: int32(id),
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
