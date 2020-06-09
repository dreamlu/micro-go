package main

import (
	"github.com/dreamlu/gt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/plugin"
	"micro-go/api-gateway/wrapper/filter"
	"micro-go/commons/wrapper/breaker"
)

// main.go
func init() {
	//token := &token.Token{}
	//token.InitConfig("127.0.0.1:8500", "micro", "config", "jwt-key", "key")

	_ = plugin.Register(
		plugin.NewPlugin(
			plugin.WithName("filter"),
			plugin.WithHandler(filter.Filter()),
		),
	)
	// 熔断限流
	_ = plugin.Register(
		plugin.NewPlugin(
			plugin.WithName("breaker"),
			plugin.WithHandler(breaker.BreakerWrapper),
		),
	)

}

func main() {

	// PrometheusBoot()
	// registry
	reg := consul.NewRegistry(
		registry.Addrs(gt.Configger().GetString("app.consul.address")),
	)
	cmd.Init(
		//micro.Name("micro-go.web.api-gateway"), // no effect
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.Registry(reg),
		//micro.Address(":"+gt.Configger().GetString("app.port")),
		//micro.WrapClient(micro_hystrix.NewClientWrapper()),
	)

}

// 访问打印出prometheus中go相关参数和含义
// 取消注释即可, main 中引用
//func PrometheusBoot(){
//	http.Handle("/metrics", promhttp.Handler())
//	// 启动web服务，监听8085端口
//	go func() {
//		err := http.ListenAndServe(":8085", nil)
//		if err != nil {
//			log.Fatal("ListenAndServe: ", err)
//		}
//	}()
//}
