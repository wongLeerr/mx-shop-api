package router

import (
	"mx-shop-api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	userGroup := router.Group("user")
	{
		userGroup.GET("/list", api.GetUserList)
		userGroup.POST("/pwd_login", api.PasswordLogin)
	}
}
