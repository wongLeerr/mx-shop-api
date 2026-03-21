package forms

type AddressReqForm struct {
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	District     string `json:"district" binding:"required"`
	Address      string `json:"address" binding:"required"`
	SignerName   string `json:"signerName" binding:"required"`
	SignerMobile string `json:"signerMobile" binding:"required"`
}
