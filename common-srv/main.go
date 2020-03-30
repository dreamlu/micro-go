// @author  dreamlu
package main

import (
	"github.com/dreamlu/gt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/consul"
	"log"
	"micro-go/common-srv/routers"
	"net/http"
)

func main() {
	// registry
	reg := consul.NewRegistry(
		registry.Addrs(gt.Configger().GetString("app.consul.address")),
	)

	// Create service
	service := web.NewService(
		web.Name("micro-go.web.common-srv"),
		web.Registry(reg),
		web.Address(":"+gt.Configger().GetString("app.port")),
	)

	_ = service.Init()

	// Create RESTful handler (using Gin)
	// Register Handler
	gin.SetMode(gin.DebugMode)
	// 路由
	router := routers.SetRouter()
	// 后台配置
	// 注释即可取消
	//back.SetBack(router)
	// 注册
	service.Handle("/", http.StripPrefix("/common-srv", router))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
