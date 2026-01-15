package router

import (
	"mx-shop-api/goods-web/api/goods"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(router *gin.RouterGroup) {
	GoodsGroup := router.Group("goods")
	{
		GoodsGroup.GET("/list", goods.GoodsList)
	}
}
