package initialize

import (
	"fmt"
	"mx-shop-api/user-interation-web/global"
	"mx-shop-api/user-interation-web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	s := zap.S()
	consulInfo := global.ServerConfig.ConsulInfo
	// 初始化地址服务连接
	addressConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserInteractionConf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		s.Fatal("【InitSrvConn】地址服务连接失败")
	}

	addressClient := proto.NewAddressClient(addressConn)
	global.Addresslient = addressClient

	// 初始化留言服务连接
	messageConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserInteractionConf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		s.Fatal("【InitSrvConn】留言服务连接失败")
	}

	messageClient := proto.NewMessageClient(messageConn)
	global.Messagelient = messageClient

	// 初始化收藏服务连接
	userFavConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserInteractionConf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		s.Fatal("【InitSrvConn】收藏服务连接失败")
	}

	userFavClient := proto.NewUserFavClient(userFavConn)
	global.UserFavClient = userFavClient
}
