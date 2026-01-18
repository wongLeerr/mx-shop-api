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
	}
}
