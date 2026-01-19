package router

import (
	"mx-shop-api/goods-web/api/goods"
	"mx-shop-api/goods-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(router *gin.RouterGroup) {
	GoodsGroup := router.Group("goods")
	{
		GoodsGroup.GET("/list", goods.GoodsList)
		GoodsGroup.POST("/create", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.CreateGoods)
		GoodsGroup.GET("/detail/:id", goods.GoodsDetail)
		GoodsGroup.DELETE("/delete/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.DeleteGoods)
		GoodsGroup.GET("/stock/:id", goods.Stock)                                                                   // 获取商品库存
		GoodsGroup.PUT("/update/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.UpdateGoods)              // 更新商品
		GoodsGroup.PATCH("/updateStatus/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.UpdateGoodStatus) // 更新商品状态
	}
}
