package initialize

import (
	"mx-shop-api/goods-web/middlewares"
	"mx-shop-api/goods-web/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// 健康检查
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})
	Router.Use(middlewares.Cors()) // 解决跨域

	APIV1Router := Router.Group("v1")
	router.InitGoodsRouter(APIV1Router)
	router.InitCategoryRouter(APIV1Router)
	router.InitBannerRouter(APIV1Router)

	return Router
}
