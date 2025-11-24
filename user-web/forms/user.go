package forms

type PasswordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required"` // 手机号应该是符合某种规范的，怎么进行校验？使用自定义 validator
	Password string `form:"password" json:"password" binding:"required,min=5,max=20"`
}
