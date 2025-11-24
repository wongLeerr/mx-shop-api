package main

import (
	"fmt"
	"mx-shop-api/user-web/global"
	"mx-shop-api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	// 初始化logger
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	// 初始化router
	Router := initialize.Routers()

	s := zap.S() // 创建sugarLogger实例
	s.Infof("🚀server will running at port: %d", global.ServerConfig.Port)
	err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		s.Panic("😭server run error:", err.Error())
	}
}
