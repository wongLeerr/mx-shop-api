package global

import (
	"mx-shop-api/user-web/config"
	"mx-shop-api/user-web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	UserSrvClient proto.UserClient
)
