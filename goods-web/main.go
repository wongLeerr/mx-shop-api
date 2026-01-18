package main

import (
	"fmt"
	"mx-shop-api/goods-web/global"
	"mx-shop-api/goods-web/initialize"
	"mx-shop-api/goods-web/utils"
	"mx-shop-api/goods-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"
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

	// 注册至Consul注册中心
	serviceId := uuid.NewV4()
	serviceIdStr := fmt.Sprintf("%s", serviceId)
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, serviceIdStr, global.ServerConfig.Name, global.ServerConfig.Tags)

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
