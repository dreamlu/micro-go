// @author  dreamlu
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/web"
	"log"
	"micro-go/common-srv/routers"
	"micro-go/commons/util/config"
	"net/http"
)

func main() {
	// registry
	registry := consul.NewRegistry(consul.Config(
		&api.Config{
			Address: config.Config.GetString("app.consul.address"),
			Scheme:  config.Config.GetString("app.consul.scheme"),
		}))

	// Create service
	service := web.NewService(
		web.Name("micro-go.web.common-srv"),
		web.Registry(registry),
		web.Address(":"+config.Config.GetString("app.port")),
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
