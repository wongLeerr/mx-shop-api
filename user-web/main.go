package main

import (
	"fmt"
	"mx-shop-api/user-web/global"
	"mx-shop-api/user-web/initialize"
	"mx-shop-api/user-web/utils"
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

	// 这里应该使用viper获取环境变量读取是否是线上环境，这里先写死
	isDebug := true
	// 开发环境希望端口号固定，不希望自动分配端口号
	if !isDebug {
		port, err := utils.GetFreeAddr()
		// err 为空，证明没报错
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

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
