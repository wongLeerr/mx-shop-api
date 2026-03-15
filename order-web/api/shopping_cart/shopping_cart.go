package shoppingcart

import (
	"context"
	"mx-shop-api/order-web/api"
	"mx-shop-api/order-web/forms"
	"mx-shop-api/order-web/global"
	"mx-shop-api/order-web/models"
	"mx-shop-api/order-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 添加商品到购物车
func Create(ctx *gin.Context) {
	s := zap.S()
	var shopcartItem forms.ShoppingCartItemForm
	err := ctx.ShouldBind(&shopcartItem)
	if err != nil {
		s.Errorln("【Create】表单验证失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 1. 校验前端传来的商品是否存在
	_, err = global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: shopcartItem.GoodsId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		s.Errorln("【Create】查询商品详情失败", err)
		return
	}

	// 2. 校验商品库存是否充足
	resp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodInvInfo{
		GoodId: shopcartItem.GoodsId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		s.Errorln("【Create】查询商品库存信息失败", err)
		return
	}

	models.ToStringLog(resp)

	if resp.Num < shopcartItem.Nums {
		api.HandleGrpcErrorToHttp(err, ctx)
		s.Errorln("【Create】商品库存不足", err)
		return
	}

	uid, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(uid.(uint)),
		GoodsId: shopcartItem.GoodsId,
		Nums:    shopcartItem.Nums,
		Checked: *shopcartItem.Checked,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		s.Errorln("【Create】添加到购物车失败", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func List(ctx *gin.Context) {
	s := zap.S()
	uid, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "缺少uid",
		})
	}
	resp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(uid.(uint)), // uid 是interface{}类型，想做类型转换必须进行断言
	})
	if err != nil {
		s.Errorln("【List】查询购物车列表失败", err)

		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	idToCartItemMap := make(map[int32]*proto.ShoppingCartInfoResponse)
	ids := make([]int32, 0)
	for _, cartItem := range resp.Data {
		ids = append(ids, cartItem.GoodsId)
		idToCartItemMap[cartItem.GoodsId] = cartItem // 维护购物车中每个条目商品id和信息的对应关系
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
			"data":  []interface{}{},
		})
		return
	}

	goodsResp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		s.Errorln("【List】批量查询商品列表失败", err)

		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	finalData := make([]interface{}, 0)
	for _, good := range goodsResp.Data {
		cartItemData := idToCartItemMap[good.Id]
		item := map[string]interface{}{
			"id":          cartItemData.Id,
			"goods_id":    good.Id,
			"goods_name":  good.Name,
			"goods_price": good.ShopPrice,
			"goods_img":   good.GoodsFrontImage,
			"nums":        cartItemData.Nums,
			"checked":     cartItemData.Checked,
		}

		finalData = append(finalData, item)
	}

	/*
		Resp:
		{
			total:12,
			data:[
				{
					id:1, // 在购物车条目中的id
					goods_id:1,
					商品信息...
					nums:数量,
					checked:是否被选中
				}
			]
		}
	*/

	ctx.JSON(http.StatusOK, gin.H{
		"total": resp.Total,
		"data":  finalData,
	})

}

func Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idInt, _ := strconv.Atoi(idStr)
	uid, _ := ctx.Get("userId")

	_, err := global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		Id:     int32(idInt),
		UserId: int32(uid.(uint)),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func Update(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	idInt, _ := strconv.Atoi(idStr)
	uid, _ := ctx.Get("userId")
	var form forms.UpdateShopCartItemForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		s.Errorln("【Update】表单验证失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(uid.(uint)),
		Id:      int32(idInt),
		Nums:    form.Nums,
		Checked: *form.Checked,
	})
	if err != nil {
		s.Errorln("【Update】更新购物车商品失败", err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}
