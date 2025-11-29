package forms

type PasswordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号应该是符合某种规范的，怎么进行校验？使用自定义 validator
	Password  string `form:"password" json:"password" binding:"required,min=5,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=4,max=4"` // 验证码
	CaptchaId string `form:"captcha_id" json:"captcha_id" bindind:"required"`       // 验证码id
}
