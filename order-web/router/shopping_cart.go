package router

import (
	shoppingcart "mx-shop-api/order-web/api/shopping_cart"
	"mx-shop-api/order-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitShopingCartRouter(router *gin.RouterGroup) {
	ShoppingCartGroup := router.Group("shopping_cart").Use(middlewares.JWTAuth())
	{
		ShoppingCartGroup.POST("/create", shoppingcart.Create)
		ShoppingCartGroup.GET("/list", shoppingcart.List)
		ShoppingCartGroup.DELETE("/:id", shoppingcart.Delete)
		ShoppingCartGroup.PATCH("/:id", shoppingcart.Update)
	}
}
