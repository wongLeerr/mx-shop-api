package global

import (
	"mx-shop-api/user-interation-web/config"
	"mx-shop-api/user-interation-web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	AddressClient proto.AddressClient
	MessageClient proto.MessageClient
	UserFavClient proto.UserFavClient
)
