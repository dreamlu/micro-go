// @author  dreamlu
package main

import (
	"demo/base-srv/routers"
	"demo/base-srv/routers/dreamlu"
	"demo/base-srv/util/cron"
	dreamlu2 "demo/base-srv/util/db/dreamlu"
	"demo/base-srv/util/pprof"
	"github.com/dreamlu/go-micro/v2/registry"
	"github.com/dreamlu/go-micro/v2/registry/consul"
	"github.com/dreamlu/go-micro/v2/web"
	"github.com/dreamlu/gt/tool/conf"
	"github.com/dreamlu/gt/tool/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// registry
	reg := consul.NewRegistry(
		registry.Addrs(conf.GetString("app.consul.address")),
	)

	// Create service
	service := web.NewService(
		web.Name("demo.api.base-srv"),
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
	log.DefaultFileLog()
	// 注册
	service.Handle("/", http.StripPrefix("/base-srv", router))
	//service.Handle("/", hystrix.BreakerWrapper(http.StripPrefix("/user-srv", router)))

	// Run server
	if err := service.Run(); err != nil {
		log.Error(err)
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
