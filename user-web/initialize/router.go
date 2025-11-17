package initialize

import (
	"mx-shop-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	APIV1Router := Router.Group("v1")
	router.InitUserRouter(APIV1Router)

	return Router
}
