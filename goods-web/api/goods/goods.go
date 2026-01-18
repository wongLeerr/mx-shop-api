package goods

import (
	"context"
	"mx-shop-api/goods-web/forms"
	"mx-shop-api/goods-web/global"
	"mx-shop-api/goods-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func GoodRespAdapter(value *proto.GoodsInfoResponse) interface{} {
	return map[string]interface{}{
		"id":          value.Id,
		"name":        value.Name,
		"goods_brief": value.GoodsBrief,
		"desc":        value.GoodsDesc,
		"ship_free":   value.ShipFree,
		"images":      value.Images,
		"desc_images": value.DescImages,
		"front_image": value.GoodsFrontImage,
		"shop_price":  value.ShopPrice,
		"category": map[string]interface{}{
			"id":   value.Category.Id,
			"name": value.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   value.Brand.Id,
			"name": value.Brand.Name,
			"logo": value.Brand.Logo,
		},
		"is_hot":  value.IsHot,
		"is_new":  value.IsNew,
		"on_sale": value.OnSale,
	}
}

func GoodsList(ctx *gin.Context) {
	s := zap.S()
	req := proto.GoodsFilterRequest{}
	// 价格最小值
	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin) // 注意这里不对err做处理，是因为即使出错了priceMinInt为0，等于没加过滤条件
	req.PriceMin = int32(priceMinInt)
	// 价格最大值
	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	req.PriceMax = int32(priceMaxInt)
	// 是否最热
	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		req.IsHot = true
	}
	// 是否最新
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		req.IsNew = true
	}
	// 是否tab栏
	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		req.IsTab = true
	}

	categoryId := ctx.DefaultQuery("c", "0")
	if categoryId != "0" {
		topCategory, _ := strconv.Atoi(categoryId)
		req.TopCategory = int32(topCategory)
	}

	page := ctx.DefaultQuery("pn", "1")
	if page != "0" {
		pageS, _ := strconv.Atoi(page)
		req.Pages = int32(pageS)
	}

	pageNum := ctx.DefaultQuery("pnum", "10")
	if pageNum != "0" {
		pageSize, _ := strconv.Atoi(pageNum)
		req.PagePerNums = int32(pageSize)
	}

	keywords := ctx.DefaultQuery("q", "")
	if keywords != "" {
		req.KeyWords = keywords
	}

	brandIdStr := ctx.DefaultQuery("b", "0")
	if brandIdStr != "0" {
		brand, _ := strconv.Atoi(brandIdStr)
		req.Brand = int32(brand)
	}

	resp, err := global.GoodSrvClient.GoodsList(context.Background(), &req)
	if err != nil {
		s.Errorf("【GoodsList】Error", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	goodsList := make([]interface{}, 0)

	for _, value := range resp.Data {
		goodsList = append(goodsList, GoodRespAdapter(value))
	}

	respMap := map[string]interface{}{
		"total": resp.Total,
		"data":  goodsList,
	}

	ctx.JSON(http.StatusOK, respMap)
}

func CreateGoods(ctx *gin.Context) {
	s := zap.S()
	var goodsForm forms.GoodsForm
	if err := ctx.ShouldBind(&goodsForm); err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := global.GoodSrvClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		CategoryId:      goodsForm.CategoryId,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		ShipFree:        *goodsForm.ShipFree,
		GoodsFrontImage: goodsForm.FrontImage,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		s.Errorf("【CreateGoods】Error", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// TODO: 库存管理
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

func GoodsDetail(ctx *gin.Context) {
	s := zap.S()
	goodId := ctx.Param("id")
	goodIdInt, _ := strconv.Atoi(goodId)
	resp, err := global.GoodSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: int32(goodIdInt),
	})
	if err != nil {
		s.Errorf("【GoodsDetail】Error", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := GoodRespAdapter(resp)

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": data,
	})
}
