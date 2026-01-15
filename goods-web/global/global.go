package global

import (
	"mx-shop-api/goods-web/config"
	"mx-shop-api/goods-web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	NacosConfig   *config.NacosConfig  = &config.NacosConfig{}
	GoodSrvClient proto.GoodsClient
)
