package main

import (
	"context"
	"fmt"
	"log"
	"net"
	stb_server "stb_consul/external_service/stb_server"
	"stb_consul/external_service/stbserver"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthImpl 健康检查实现
type HealthImpl struct{}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

//Watch 这个没用，只是为了让HealthImpl实现RegisterHealthServer内部的interface接口
func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

//grpc开启
func externalServer() {
	lis, err := net.Listen("tcp", ":3001")
	if err != nil {
		logrus.Info("外置服务开启失败:", err)
		panic(err)
	}
	logrus.WithFields(logrus.Fields{
		"tcp": ":3001",
	}).Info("external server")
	s := grpc.NewServer()
	stbserver.RegisterStbServerServer(s, &stb_server.StbServe{})
	grpc_health_v1.RegisterHealthServer(s, &HealthImpl{})
	s.Serve(lis)
	log.Println("grpc start")
}

//grpc注册进consul
func grpcRegister() {
	config := api.DefaultConfig()
	config.Address = consulAddress
	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}
	agent := client.Agent()

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", "StbServe", localIP, localPort), // 服务节点的名称
		Name:    fmt.Sprintf("grpc.health.v1.%v", "StbServe"),            // 服务名称
		Tags:    []string{"StbServe"},                                    // tag，可以为空
		Port:    localPort,                                               // 服务端口
		Address: localIP,                                                 // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval: "5s", // 健康检查间隔
			// grpc 支持，执行健康检查的地址，service 会传到 Health.Check 函数中
			GRPC:                           fmt.Sprintf("%v:%v/%v", localIP, localPort, "StbServe"),
			DeregisterCriticalServiceAfter: "5s", // 注销时间，相当于过期时间
		},
	}

	if err := agent.ServiceRegister(reg); err != nil {
		panic(err)
	}
}

func grpcLoad() {
	grpcRegister()
	externalServer()
}
