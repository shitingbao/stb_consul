package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

func client() {
	var lastIndex uint64
	config := api.DefaultConfig()
	config.Address = "124.70.156.31:8500" //consul server

	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println("api new client is failed, err:", err)
		return
	}
	services, metainfo, err := client.Health().Service("service_shitingbao", "shitingbao_test_service", true, &api.QueryOptions{
		WaitIndex: lastIndex, // 同步点，这个调用将一直阻塞，直到有新的更新
	})
	if err != nil {
		logrus.Panic("error retrieving instances from Consul:", err)
	}
	lastIndex = metainfo.LastIndex

	for _, service := range services {
		fmt.Println("service.Service.Address:", service.Service.Address, "service.Service.Port:", service.Service.Port)
		fmt.Println("service.Service:", service.Service)
	}
}
