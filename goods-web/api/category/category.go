package category

import (
	"context"
	"encoding/json"
	"mx-shop-api/goods-web/api"
	"mx-shop-api/goods-web/forms"
	"mx-shop-api/goods-web/global"
	"mx-shop-api/goods-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func CategoryList(ctx *gin.Context) {
	s := zap.S()
	resp, err := global.GoodSrvClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		s.Errorf("【CategoryList】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	respSlice := make([]interface{}, 0)
	err = json.Unmarshal([]byte(resp.JsonData), &respSlice)
	if err != nil {
		s.Errorf("反序列化失败：", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": respSlice,
	})
}

func CreateCategory(ctx *gin.Context) {
	s := zap.S()
	var createCategoryForm forms.CategoryInfoForm
	err := ctx.ShouldBind(&createCategoryForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := global.GoodSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           createCategoryForm.Name,
		ParentCategory: createCategoryForm.ParentCategoryId,
		Level:          createCategoryForm.Level,
		IsTab:          *createCategoryForm.IsTab,
	})
	if err != nil {
		s.Errorf("【CreateCategory】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": resp,
	})
}

func DeleteCategory(ctx *gin.Context) {
	s := zap.S()
	categoryIdStr := ctx.Param("id")
	categoryId, _ := strconv.Atoi(categoryIdStr)
	_, err := global.GoodSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{
		Id: int32(categoryId),
	})
	if err != nil {
		s.Errorf("【DeleteCategory】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func UpdateCategory(ctx *gin.Context) {
	s := zap.S()
	var updateCategoryForm forms.UpdateCategoryForm
	err := ctx.ShouldBind(&updateCategoryForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	categoryIdStr := ctx.Param("id")
	categoryId, _ := strconv.Atoi(categoryIdStr)
	_, err = global.GoodSrvClient.UpdateCategory(context.Background(), &proto.CategoryInfoRequest{
		Id:    int32(categoryId),
		Name:  updateCategoryForm.Name,
		IsTab: *updateCategoryForm.IsTab,
	})
	if err != nil {
		s.Errorf("【UpdateCategory】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func SubCategory(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.DefaultQuery("id", "0")
	categoryId, _ := strconv.Atoi(idStr)
	resp, err := global.GoodSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(categoryId),
	})
	if err != nil {
		s.Errorf("【SubCategory】Error", err.Error())
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	// todo：转换成前端需要的下划线格式
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": resp,
	})
}
