package brands

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
)

func CreateBrand(ctx *gin.Context) {
	s := zap.S()
	var brand forms.BrandForm
	err := ctx.ShouldBind(&brand)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := global.GoodSrvClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brand.Name,
		Logo: brand.Logo,
	})
	if err != nil {
		s.Errorf("【CreateBrand】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": resp,
	})
}

func DeleteBrand(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	_, err := global.GoodSrvClient.DeleteBrand(context.Background(), &proto.BrandRequest{
		Id: int32(id),
	})
	if err != nil {
		s.Errorf("【DeleteBrand】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func UpdateBrand(ctx *gin.Context) {
	s := zap.S()
	var brand forms.BrandForm
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	err := ctx.ShouldBind(&brand)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err = global.GoodSrvClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   int32(id),
		Name: brand.Name,
		Logo: brand.Logo,
	})
	if err != nil {
		s.Errorf("【UpdateBrand】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})

}

func BrandList(ctx *gin.Context) {
	s := zap.S()
	pnStr := ctx.DefaultQuery("pn", "1")
	pn, _ := strconv.Atoi(pnStr)
	pSizeStr := ctx.DefaultQuery("pSize", "10")
	pSize, _ := strconv.Atoi(pSizeStr)
	resp, err := global.GoodSrvClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       int32(pn),
		PagePerNums: int32(pSize),
	})
	if err != nil {
		s.Errorf("【BrandList】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
