package router

import (
	"mx-shop-api/goods-web/api/brands"

	"github.com/gin-gonic/gin"
)

func InitBrandRouter(router *gin.RouterGroup) {
	BrandGroup := router.Group("brand")
	{
		BrandGroup.POST("/create", brands.CreateBrand)
		BrandGroup.DELETE("/delete/:id", brands.DeleteBrand)
		BrandGroup.PUT("/update/:id", brands.UpdateBrand)
		BrandGroup.GET("/list", brands.BrandList)
	}
}
