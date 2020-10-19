// @author  dreamlu
package main

import (
	"github.com/dreamlu/go-micro/v2/registry"
	"github.com/dreamlu/go-micro/v2/registry/consul"
	"github.com/dreamlu/go-micro/v2/web"
	"github.com/dreamlu/gt"
	"github.com/gin-gonic/gin"
	"log"
	"micro-go/base-srv/routers"
	"micro-go/base-srv/routers/dreamlu"
	"micro-go/base-srv/util/cron"
	dreamlu2 "micro-go/base-srv/util/db/dreamlu"
	"micro-go/base-srv/util/pprof"
	"net/http"
)

func main() {
	// registry
	reg := consul.NewRegistry(
		registry.Addrs(gt.Configger().GetString("app.consul.address")),
	)

	// Create service
	service := web.NewService(
		web.Name("micro-go.api.base-srv"),
		web.Registry(reg),

		//web.Address(":"+gt.Configger().GetString("app.port")),
	)

	_ = service.Init()

	// Create RESTful handler (using Gin)
	// Register Handler
	gin.SetMode(gin.DebugMode)
	// 路由
	router := routers.Router
	pprof.Register(router)
	// out log to file
	gt.Logger().DefaultFileLog()
	// 注册
	service.Handle("/", http.StripPrefix("/base-srv", router))
	//service.Handle("/", hystrix.BreakerWrapper(http.StripPrefix("/user-srv", router)))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// 初始化
func init() {
	// 数据库模型自动生成
	dreamlu.InitRouter()
	dreamlu2.InitDB()
	// 定时任务
	go cron.Cron()
}
