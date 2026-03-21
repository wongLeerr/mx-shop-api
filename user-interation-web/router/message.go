package router

import (
	"mx-shop-api/user-interation-web/api/message"
	"mx-shop-api/user-interation-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitMessageRouter(router *gin.RouterGroup) {
	messageGroup := router.Group("message").Use(middlewares.JWTAuth())
	{
		messageGroup.POST("", message.Create)
		messageGroup.GET("/list", message.List)
	}
}
