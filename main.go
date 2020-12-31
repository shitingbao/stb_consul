package main

import (
	"fmt"

	"net/http"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	consulAddress = "124.70.156.31:8500"
	localIP       = "124.70.156.31"
	localPort     = 3001
)

func consulRegister() {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "shitingbao"
	registration.Name = "service_shitingbao"
	registration.Port = localPort
	registration.Tags = []string{"shitingbao_test_service"}
	registration.Address = localIP

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
	check.Timeout = "5s"
	check.Interval = "5s"
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check

	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
}

//Handler 3001
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you are visiting health check api:3001"))
}

func main() {
	consulRegister()
	//定义一个http接口
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
}
