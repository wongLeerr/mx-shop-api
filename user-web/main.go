package main

import (
	"fmt"
	"mx-shop-api/user-web/global"
	"mx-shop-api/user-web/initialize"
	customValidator "mx-shop-api/user-web/validator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	// 初始化logger
	initialize.InitLogger()
	// 初始化配置文件
	initialize.InitConfig()
	// 初始化router
	Router := initialize.Routers()
	// 初始化srv的连接，生成全局client
	initialize.InitSrvConn()
	// 注册自定义表单验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", customValidator.ValidateMobile)
	}

	s := zap.S() // 创建sugarLogger实例
	s.Infof("🚀server will running at port: %d", global.ServerConfig.Port)
	err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		s.Panic("😭server run error:", err.Error())
	}
}
