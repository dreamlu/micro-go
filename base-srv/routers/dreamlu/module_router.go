package dreamlu

import (
	"demo/base-srv/controllers/module/day_sign"
	"demo/base-srv/controllers/module/day_sign/img"
	"demo/base-srv/controllers/module/day_sign/text"
	"demo/base-srv/controllers/module/live_bro"
	"demo/base-srv/routers"
)

func ModuleRouter() {
	//组的路由,version
	v := routers.V
	// 模块--小程序直播
	live_bros := v.Group("/live_bro")
	{
		live_bros.POST("/list", live_bro.LiveBroList)
	}

	// 模块--日签
	day_signs := v.Group("/day_sign")
	{
		day_signs.GET("/search", day_sign.GetBySearch)
		day_signs.GET("/id", day_sign.GetByID)
		day_signs.DELETE("/delete/:id", day_sign.Delete)
		day_signs.POST("/create", day_sign.Create)
		day_signs.PUT("/update", day_sign.Update)

		// 日签--图片库
		imgs := day_signs.Group("/img")
		{
			imgs.GET("/search", img.GetBySearch)
			imgs.GET("/id", img.GetByID)
			imgs.DELETE("/delete/:id", img.Delete)
			imgs.POST("/create", img.Create)
			imgs.PUT("/update", img.Update)
		}
		// 日签--文字库
		texts := day_signs.Group("/text")
		{
			texts.GET("/search", text.GetBySearch)
			texts.GET("/id", text.GetByID)
			texts.DELETE("/delete/:id", text.Delete)
			texts.POST("/create", text.Create)
			texts.PUT("/update", text.Update)
		}
	}
}
