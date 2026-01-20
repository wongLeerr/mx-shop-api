package forms

type BannerInfo struct {
	Image string `form:"image" json:"image" binding:"required"`
	Url   string `form:"url" json:"url" binding:"required"`
	Index int32  `form:"index" json:"index"`
}
