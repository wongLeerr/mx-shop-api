package router

import (
	"mx-shop-api/goods-web/api/category"

	"github.com/gin-gonic/gin"
)

func InitCategoryRouter(router *gin.RouterGroup) {
	GoodsGroup := router.Group("category")
	{
		GoodsGroup.GET("/list", category.CategoryList)
		GoodsGroup.POST("/create", category.CreateCategory) //  middlewares.JWTAuth(), middlewares.IsAdmin(),
		GoodsGroup.DELETE("/delete/:id", category.DeleteCategory)
		GoodsGroup.PUT("/update/:id", category.UpdateCategory)
		GoodsGroup.GET("/subCategory", category.SubCategory)
	}
}
