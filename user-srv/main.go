// @author  dreamlu
package main

import (
	der "github.com/dreamlu/go-tool"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/web"
	"log"
	"micro-go/user-srv/routers"
	"net/http"
)

func main() {

	// registry
	registry := consul.NewRegistry(consul.Config(
		&api.Config{
			Address: der.Configger().GetString("app.consul.address"),
			Scheme:  der.Configger().GetString("app.consul.scheme"),
		}))

	// Create service
	service := web.NewService(
		web.Name("micro-go.web.user-srv"),
		web.Registry(registry),
		web.Address(":"+der.Configger().GetString("app.port")),
	)

	_ = service.Init()

	// Create RESTful handler (using Gin)
	// Register Handler
	gin.SetMode(gin.DebugMode)
	// 路由
	router := routers.SetRouter()
	// out log to file
	der.Logger().DefaultFileLog()
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
