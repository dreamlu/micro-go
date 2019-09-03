package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/micro/micro/cmd"
	"github.com/micro/micro/plugin"
	"micro-go/api-gateway/wrapper/filter"
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
}

func main() {

	// PrometheusBoot()
	cmd.Init(micro.WrapHandler(prometheus.NewHandlerWrapper()))
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
