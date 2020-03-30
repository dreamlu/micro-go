// @author  dreamlu
package main

import (
	"github.com/dreamlu/gt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"micro-go/user-srv/routers"
	"net/http"
)

func main() {
	// registry
	reg := consul.NewRegistry(
		registry.Addrs(gt.Configger().GetString("app.consul.address")),
	)

	// Create service
	service := web.NewService(
		web.Name("micro-go.web.user-srv"),
		web.Registry(reg),
		web.Address(":"+gt.Configger().GetString("app.port")),
	)

	_ = service.Init()

	// Create RESTful handler (using Gin)
	// Register Handler
	gin.SetMode(gin.DebugMode)
	// 路由
	router := routers.SetRouter()
	// out log to file
	gt.Logger().DefaultFileLog()
	// 注册
	service.Handle("/", http.StripPrefix("/user-srv", router))
	//service.Handle("/", hystrix.BreakerWrapper(http.StripPrefix("/user-srv", router)))

	//c := client.NewClient(
	//	client.Selector(
	//		selector.NewSelector(
	//			selector.Registry(memory.NewRegistry()),
	//		),
	//	),
	//)
	//req := c.NewRequest("micro-go.web.common-srv", "basic.basic",
	//	map[string]string{
	//		"foo": "bar",
	//	}, client.WithContentType("application/json"))
	//var rsp map[string]interface{}
	//c.Call(context.TODO(), req, rsp)
	//log.Println(rsp)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
