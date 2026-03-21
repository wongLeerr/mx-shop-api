package router

import (
	"mx-shop-api/user-interation-web/api/address"
	"mx-shop-api/user-interation-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitAddressRouter(router *gin.RouterGroup) {
	addressGroup := router.Group("address").Use(middlewares.JWTAuth())
	{
		addressGroup.POST("", address.Create)
		addressGroup.DELETE("/:id", address.Del)
		addressGroup.PATCH("/:id", address.Update)
		addressGroup.GET("/list", address.List)
	}
}
