package global

import (
	"mx-shop-api/user-interation-web/config"
	"mx-shop-api/user-interation-web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	Addresslient  proto.AddressClient
	Messagelient  proto.MessageClient
	UserFavClient proto.UserFavClient
)
