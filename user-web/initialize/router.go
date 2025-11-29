package initialize

import (
	"mx-shop-api/user-web/middlewares"
	"mx-shop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors()) // 解决跨域

	APIV1Router := Router.Group("v1")
	router.InitUserRouter(APIV1Router)

	return Router
}
