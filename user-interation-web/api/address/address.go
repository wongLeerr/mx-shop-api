package address

import (
	"context"
	"mx-shop-api/user-interation-web/forms"
	"mx-shop-api/user-interation-web/global"
	"mx-shop-api/user-interation-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Create(ctx *gin.Context) {
	s := zap.S()
	uid, _ := ctx.Get("userId")
	var addressForm forms.AddressReqForm
	err := ctx.ShouldBind(&addressForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := global.AddressClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       int32(uid.(uint)),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
		"id":  resp.Id,
	})
}

func Del(ctx *gin.Context) {
	s := zap.S()
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	uid, _ := ctx.Get("userId")

	_, err := global.AddressClient.DeleteAddress(context.Background(), &proto.AddressRequest{
		Id:     int32(id),
		UserId: int32(uid.(uint)),
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func Update(ctx *gin.Context) {
	s := zap.S()
	uid, _ := ctx.Get("userId")
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	var addressForm forms.AddressReqForm
	err := ctx.ShouldBind(&addressForm)
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	_, err = global.AddressClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           int32(id),
		UserId:       int32(uid.(uint)),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func List(ctx *gin.Context) {
	s := zap.S()
	uid, _ := ctx.Get("userId")
	rsp, err := global.AddressClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: int32(uid.(uint)),
	})
	if err != nil {
		s.Errorln(err.Error())
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	resp := make(map[string]interface{})
	data := make([]interface{}, 0)
	resp["total"] = rsp.Total
	for _, value := range rsp.Data {
		data = append(data, map[string]interface{}{
			"id":           value.Id,
			"userId":       value.UserId,
			"province":     value.Province,
			"city":         value.City,
			"district":     value.District,
			"address":      value.Address,
			"signerName":   value.SignerName,
			"signerMobile": value.SignerMobile,
		})
	}
	resp["data"] = data
	/**
	rsp:
	{
		total:1,
		data:[
			{
				id:1,
				userId:1,
				province:"",
				city:"",
				district:"",
				address:"",
				signerName:"",
				signerMobile:"",
			}
		]
	}
	*/

	ctx.JSON(http.StatusOK, rsp)
}
