package db

import (
	der "github.com/dreamlu/go-tool"
)

// init param
var DBTool = &der.DBTool{}

//var Log = &der.Log{}

func init() {
	DBTool.NewDBTool()
	DBTool.DB.LogMode(true)
	//Log.DefaultFileLog()
}
