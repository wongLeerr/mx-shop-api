package router

import (
	"mx-shop-api/order-web/api/order"
	"mx-shop-api/order-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(router *gin.RouterGroup) {
	OrderGroup := router.Group("order").Use(middlewares.JWTAuth()) // 对该组别下的所有API请求都进行登录鉴权验证
	{
		OrderGroup.POST("", order.Create)
		OrderGroup.GET("/:id", order.Detail)
		OrderGroup.GET("/list", order.List)
	}
}
