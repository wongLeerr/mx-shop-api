package initialize

import (
	"mx-shop-api/goods-web/middlewares"
	"mx-shop-api/goods-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors()) // 解决跨域

	APIV1Router := Router.Group("v1")
	router.InitGoodsRouter(APIV1Router)

	return Router
}
