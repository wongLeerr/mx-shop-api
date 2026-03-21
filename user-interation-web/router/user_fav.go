package router

import (
	userfav "mx-shop-api/user-interation-web/api/user_fav"
	"mx-shop-api/user-interation-web/middlewares"

	"github.com/gin-gonic/gin"
)

func InitUserFavRouter(router *gin.RouterGroup) {
	userFavGroup := router.Group("user_fav").Use(middlewares.JWTAuth())
	{
		userFavGroup.POST("", userfav.Create)
		userFavGroup.DELETE("/:id", userfav.Del)
		userFavGroup.GET("/list", userfav.List)
		userFavGroup.GET("/:id", userfav.Get)
	}
}
