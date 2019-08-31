package main

import (
	"github.com/micro/go-micro"
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

const name = "API gateway"

func main() {

	cmd.Init(
		micro.WrapClient(),
		)
}
