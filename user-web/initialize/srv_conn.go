package initialize

import (
	"fmt"
	"mx-shop-api/user-web/global"
	"mx-shop-api/user-web/proto"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	s := zap.S()
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvConf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		s.Fatal("【InitSrvConn】用户服务连接失败")
	}

	userClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userClient
}

// 未使用负载均衡版本（直接call一个确定服务的srv）
func InitSrvCon2() {
	s := zap.S()
	// 从注册中心获取用户服务（服务发现）
	conf := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	conf.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	client, err := api.NewClient(conf)
	if err != nil {
		s.Errorw("gen new client:", err.Error())
	}

	s.Infof("🐶🐶🐶 %s", global.ServerConfig.UserSrvConf.Name)
	service, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvConf.Name))
	if err != nil {
		s.Errorw("get service err:", err.Error())
	}
	var userSrvHost string
	var userSrvPort int
	for _, value := range service {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}

	// 拨号连接user grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		s.Errorw("connect to user service error:", err.Error())
		return
	}

	// 生成grpc的client并调用接口
	userClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userClient
}
