package config

import (
	der "github.com/dreamlu/go-tool"
)

// init param
var Config = &der.Config{}

//var Log = &der.Log{}

func init() {
	Config.NewConfig()
	//Log.DefaultFileLog()
}
