package initialize

import (
	"fmt"
	"mx-shop-api/order-web/global"
	"mx-shop-api/order-web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	s := zap.S()
	consulInfo := global.ServerConfig.ConsulInfo
	// 初始化订单底层服务连接
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.OrderSrvConf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		s.Fatal("【InitSrvConn】订单服务连接失败")
	}

	orderClient := proto.NewOrderClient(orderConn)
	global.OrderSrvClient = orderClient

	// 初始化商品底层服务连接
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvConf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		s.Fatal("【InitSrvConn】商品服务连接失败")
	}

	goodsClient := proto.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsClient

	// 初始化库存底层服务连接
	inventoryConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.InventorySrvConf.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		s.Fatal("【InitSrvConn】库存服务连接失败")
	}

	inventoryClient := proto.NewInventoryClient(inventoryConn)
	global.InventorySrvClient = inventoryClient
}
