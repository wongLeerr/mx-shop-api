package main

import (
	"fmt"
	"mx-shop-api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	// 初始化logger
	initialize.InitLogger()
	// 初始化router
	Router := initialize.Routers()

	PORT := 8021
	s := zap.S() // 创建sugarLogger实例
	s.Infof("🚀server will running at port: %d", PORT)
	err := Router.Run(fmt.Sprintf(":%d", PORT))
	if err != nil {
		s.Panic("😭server run error:", err.Error())
	}
}
