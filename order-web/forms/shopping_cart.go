package forms

type ShoppingCartItemForm struct {
	GoodsId int32 `json:"goods_id" form:"goods_id" binding:"required"`
	Checked *bool `json:"checked" form:"checked"`
	Nums    int32 `json:"nums" form:"nums" binding:"required,min=1"`
}

type UpdateShopCartItemForm struct {
	Checked *bool `json:"checked" form:"checked" binding:"required"`
	Nums    int32 `json:"nums" form:"nums" binding:"required,min=1"`
}
