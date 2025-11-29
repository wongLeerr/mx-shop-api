package router

import (
	"mx-shop-api/user-web/api"
	"mx-shop-api/user-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.RouterGroup) {
	userGroup := router.Group("user")
	{
		userGroup.GET("/list", middlewares.JWTAuth(), middlewares.IsAdmin(), api.GetUserList) // 鉴权是否登录、鉴权是否是管理员、前两者都满足才会返回userlist
		userGroup.POST("/pwd_login", api.PasswordLogin)
	}
}
