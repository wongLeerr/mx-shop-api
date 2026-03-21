package forms

type AddressReqForm struct {
	Province     string `json:"province" binding:"required"`
	City         string `json:"city" binding:"required"`
	District     string `json:"district" binding:"required"`
	Address      string `json:"address" binding:"required"`
	SignerName   string `json:"signerName" binding:"required"`
	SignerMobile string `json:"signerMobile" binding:"required"`
}

type MessageReqForm struct {
	MessageType int32  `json:"messageType" binding:"required,oneof=1 2 3 4 5"`
	Subject     string `json:"subject" binding:"required"`
	Message     string `json:"message"`
	File        string `json:"file"`
}
