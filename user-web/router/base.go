package router

import (
	"mx-shop-api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(router *gin.RouterGroup) {
	baseGroup := router.Group("base")
	{
		baseGroup.GET("captcha", api.GetCaptcha)
	}
}
