// 本文件为商品分类和品牌关联关系api
package brands

import (
	"context"
	"mx-shop-api/goods-web/api"
	"mx-shop-api/goods-web/forms"
	"mx-shop-api/goods-web/global"
	"mx-shop-api/goods-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 建立品牌和分类的关联关系
func CreateCategoryBrand(ctx *gin.Context) {
	s := zap.S()
	var categoryBrandForm forms.CategoryBrandForm
	err := ctx.ShouldBind(&categoryBrandForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := global.GoodSrvClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: categoryBrandForm.CategoryId,
		BrandId:    categoryBrandForm.BrandId,
	})
	if err != nil {
		s.Errorf("【CreateCategoryBrand】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	respData := make(map[string]interface{})
	respData["id"] = resp.Id

	ctx.JSON(http.StatusOK, respData)
}

// 删除品牌和分类关联关系
func DeleteCategoryBrand(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err = global.GoodSrvClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id: int32(id),
	})
	if err != nil {
		s.Errorf("【DeleteCategoryBrand】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

// 更改品牌和分类关联关系
func UpdateCategoryBrand(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var categoryBrandForm forms.CategoryBrandForm
	err = ctx.ShouldBind(&categoryBrandForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	_, err = global.GoodSrvClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         int32(id),
		CategoryId: categoryBrandForm.CategoryId,
		BrandId:    categoryBrandForm.BrandId,
	})
	if err != nil {
		s.Errorf("【UpdateCategoryBrand】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

// 获取某个分类下的所有品牌
func GetBrandByCategoryId(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.DefaultQuery("id", "0")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := global.GoodSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: int32(id),
	})
	if err != nil {
		s.Errorf("【GetBrandByCategoryId】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func CategoryBrandList(ctx *gin.Context) {
	s := zap.S()
	pnStr := ctx.DefaultQuery("pn", "1")
	pn, _ := strconv.Atoi(pnStr)
	pSizeStr := ctx.DefaultQuery("pSize", "10")
	pSize, _ := strconv.Atoi(pSizeStr)

	response, err := global.GoodSrvClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
		Pages:       int32(pn),
		PagePerNums: int32(pSize),
	})
	if err != nil {
		s.Errorf("【CategoryBrandList】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	Resp := make(map[string]interface{}, 0)
	data := make([]interface{}, 0)
	for _, resp := range response.Data {
		respData := make(map[string]interface{})
		respData["id"] = resp.Id
		categoryData := make(map[string]interface{})
		categoryData["id"] = resp.Category.Id
		categoryData["name"] = resp.Category.Name
		categoryData["level"] = resp.Category.Level
		categoryData["is_tab"] = resp.Category.IsTab
		categoryData["parent"] = resp.Category.ParentCategory
		respData["category"] = categoryData
		brandData := make(map[string]interface{})
		brandData["id"] = resp.Brand.Id
		brandData["name"] = resp.Brand.Name
		brandData["logo"] = resp.Brand.Logo
		respData["brand"] = brandData
		data = append(data, respData)
	}
	Resp["total"] = response.Total
	Resp["data"] = data
	ctx.JSON(http.StatusOK, Resp)
}
