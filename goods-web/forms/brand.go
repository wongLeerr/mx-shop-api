package forms

type BrandForm struct {
	Name string `form:"name" json:"name" binding:"required,min=1,max=20"`
	Logo string `form:"logo" json:"logo" binding:"required"`
}
