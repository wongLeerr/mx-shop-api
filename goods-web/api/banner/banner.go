package banner

import (
	"context"
	"mx-shop-api/goods-web/api"
	"mx-shop-api/goods-web/forms"
	"mx-shop-api/goods-web/global"
	"mx-shop-api/goods-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func CreateBanner(ctx *gin.Context) {
	s := zap.S()
	var bannerInfo forms.BannerInfo
	err := ctx.ShouldBind(&bannerInfo)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := global.GoodSrvClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Image: bannerInfo.Image,
		Url:   bannerInfo.Url,
		Index: bannerInfo.Index,
	})
	if err != nil {
		s.Errorf("【CreateBanner】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": resp,
	})
}

func DeleteBanner(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)

	_, err := global.GoodSrvClient.DeleteBanner(context.Background(), &proto.BannerRequest{
		Id: int32(id),
	})
	if err != nil {
		s.Errorf("【DeleteBanner】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func UpdateBanner(ctx *gin.Context) {
	s := zap.S()
	var bannerInfo forms.BannerInfo
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	err := ctx.ShouldBind(&bannerInfo)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err = global.GoodSrvClient.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(id),
		Image: bannerInfo.Image,
		Url:   bannerInfo.Url,
		Index: bannerInfo.Index,
	})
	if err != nil {
		s.Errorf("【UpdateBanner】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func BannerList(ctx *gin.Context) {
	s := zap.S()
	resp, err := global.GoodSrvClient.BannerList(context.Background(), &emptypb.Empty{})
	if err != nil {
		s.Errorf("【BannerList】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
