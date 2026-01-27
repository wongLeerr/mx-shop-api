package main

import (
	"fmt"
	"mx-shop-api/oss-web/global"
	"mx-shop-api/oss-web/initialize"
	"mx-shop-api/oss-web/utils"
	"mx-shop-api/oss-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func main() {
	//1. 初始化logger
	initialize.InitLogger()

	//2. 初始化配置文件
	initialize.InitConfig()

	//3. 初始化routers
	Router := initialize.Routers()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

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

	serviceId := uuid.NewV4()
	serviceIdStr := fmt.Sprintf("%s", serviceId)
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceIdStr)

	s := zap.S() // 创建sugarLogger实例
	s.Infof("🚀server will running at port: %d", global.ServerConfig.Port)
	go func() {
		err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
		if err != nil {
			s.Panic("😭server run error:", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err := registerClient.DeRegister(serviceIdStr)
	if err != nil {
		s.Errorf("注销失败")
	} else {
		s.Info("注销成功")
	}
}
