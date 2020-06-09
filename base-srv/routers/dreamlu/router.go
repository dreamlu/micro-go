package dreamlu

import (
	"micro-go/base-srv/controllers/admin"
	"micro-go/base-srv/controllers/admin/applet"
	"micro-go/base-srv/controllers/admin/setup/freight"
	"micro-go/base-srv/controllers/admin/setup/freight/freight_city"
	"micro-go/base-srv/controllers/admin/setup/logistic"
	"micro-go/base-srv/controllers/ali/sms"
	"micro-go/base-srv/controllers/carousel"
	"micro-go/base-srv/controllers/client"
	"micro-go/base-srv/controllers/client/address"
	"micro-go/base-srv/controllers/client/cart"
	"micro-go/base-srv/controllers/goods"
	"micro-go/base-srv/controllers/goods/category"
	"micro-go/base-srv/controllers/goods/notice"
	"micro-go/base-srv/controllers/goods/regroup"
	"micro-go/base-srv/controllers/goods/warehouse"
	"micro-go/base-srv/controllers/module"
	"micro-go/base-srv/controllers/order"
	"micro-go/base-srv/controllers/order/comment"
	"micro-go/base-srv/controllers/order/order_refund"
	"micro-go/base-srv/controllers/order/order_refund/info"
	"micro-go/base-srv/controllers/pca"
	"micro-go/base-srv/controllers/record"
	"micro-go/base-srv/controllers/wx"
	"micro-go/base-srv/routers"
	"micro-go/commons/controllers/delivery"
)

func InitRouter() {

	//组的路由,version
	v := routers.V
	{
		//用户
		clients := v.Group("/client")
		{
			clients.GET("/token", client.Token)
			clients.GET("/search", client.GetBySearch)
			clients.GET("/idc", client.GetByIDC)
			clients.GET("/id", client.GetByID)
			clients.DELETE("/delete/:id", client.Delete)
			clients.POST("/create", client.Create)
			clients.PUT("/update", client.Update)
			//用户
			carts := clients.Group("/cart")
			{
				carts.GET("/search", cart.GetBySearch)
				carts.GET("/id", cart.GetByID)
				carts.DELETE("/delete/:id", cart.Delete)
				carts.POST("/create", cart.Create)
				carts.PUT("/update", cart.Update)
			}
			// 客户地址
			addresss := clients.Group("/address")
			{
				addresss.GET("/search", address.Search)
				addresss.GET("/id", address.Get)
				addresss.DELETE("/delete/:id", address.Delete)
				addresss.POST("/create", address.Create)
				addresss.POST("/create_more", address.CreateMore)
				addresss.PUT("/update", address.Update)
			}
		}
		// 商品
		goodss := v.Group("/goods")
		{
			goodss.GET("/search", goods.GetBySearch)
			goodss.GET("/id", goods.GetByID)
			goodss.DELETE("/delete/:id", goods.Delete)
			goodss.POST("/create", goods.Create)
			goodss.PUT("/update", goods.Update)

			// 分类
			categorys := goodss.Group("/category")
			{
				categorys.GET("/list", category.List)
				categorys.GET("/search", category.GetBySearch)
				categorys.GET("/id", category.GetByID)
				categorys.DELETE("/delete/:id", category.Delete)
				categorys.POST("/create", category.Create)
				categorys.PUT("/update", category.Update)
			}

			// 商品--复团
			regroups := goodss.Group("/regroup")
			{
				regroups.GET("/count", regroup.Count)
				regroups.GET("/search", regroup.GetBySearch)
				regroups.GET("/id", regroup.GetByID)
				regroups.DELETE("/delete/:id", regroup.Delete)
				regroups.POST("/create", regroup.Create)
				regroups.PUT("/update", regroup.Update)
			}

			// 开团提醒
			notices := goodss.Group("/notice")
			{
				notices.POST("/create", notice.Create)
			}

			// 商品--货仓
			warehouses := goodss.Group("/warehouse")
			{
				warehouses.GET("/search", warehouse.GetBySearch)
				warehouses.GET("/id", warehouse.GetByID)
				warehouses.DELETE("/delete/:id", warehouse.Delete)
				warehouses.POST("/create", warehouse.Create)
				warehouses.PUT("/update", warehouse.Update)
			}
		}
		// 轮播
		carousels := v.Group("/carousel")
		{
			carousels.GET("/search", carousel.GetBySearch)
			carousels.GET("/id", carousel.GetByID)
			carousels.DELETE("/delete/:id", carousel.Delete)
			carousels.POST("/create", carousel.Create)
			carousels.PUT("/update", carousel.Update)
		}
		// 订单
		orders := v.Group("/order")
		{
			orders.POST("/ship", order.ShipExcel)
			orders.GET("/export", order.ExportExcel)
			orders.GET("/search", order.GetBySearch)
			orders.GET("/id", order.GetByID)
			orders.DELETE("/delete/:id", order.Delete)
			orders.POST("/create", order.Create)
			orders.PUT("/update", order.Update)
			//orders.PUT("/update_more", order.UpdateMore)

			// 评论
			comments := orders.Group("/comment")
			{
				comments.GET("/search", comment.GetBySearch)
				comments.GET("/id", comment.GetByID)
				comments.DELETE("/delete/:id", comment.Delete)
				comments.POST("/create", comment.Create)
				comments.POST("/create_more", comment.CreateMore)
				comments.PUT("/update", comment.Update)
			}

			// 订单退款详情
			order_refunds := orders.Group("/refund")
			{
				order_refunds.GET("/search", order_refund.GetBySearch)
				order_refunds.GET("/id", order_refund.GetByID)
				order_refunds.DELETE("/delete/:id", order_refund.Delete)
				order_refunds.POST("/create", order_refund.Create)
				order_refunds.PUT("/update", order_refund.Update)

				order_r_infos := order_refunds.Group("/info")
				{
					order_r_infos.GET("/search", info.Search)
					order_r_infos.GET("/id", info.Get)
					order_r_infos.DELETE("/delete/:id", info.Delete)
					order_r_infos.POST("/create", info.Create)
					order_r_infos.PUT("/update", info.Update)
				}
			}
			// 快递查询
			deliverys := orders.Group("/delivery")
			{
				deliverys.GET("/get", delivery.Get)
			}
		}

		// admin
		admins := v.Group("/admin")
		{
			admins.GET("/search", admin.GetBySearch)
			admins.GET("/id", admin.GetByID)
			admins.DELETE("/delete/:id", admin.Delete)
			admins.POST("/create", admin.Create)
			admins.PUT("/update", admin.Update)
			admins.POST("/login", admin.Login)

			// applet
			applets := admins.Group("/applet")
			{
				applets.GET("/search", applet.GetBySearch)
				applets.GET("/id", applet.GetByID)
				applets.DELETE("/delete/:id", applet.Delete)
				applets.POST("/create", applet.Create)
				applets.PUT("/update", applet.Update)
				applets.POST("/download", applet.DownLoad)
			}

			// 设置
			setup := admins.Group("/setup")
			// 物流
			logistics := setup.Group("/logistic")
			{
				logistics.GET("/search", logistic.GetBySearch)
				logistics.GET("/id", logistic.GetByID)
				logistics.DELETE("/delete/:id", logistic.Delete)
				logistics.POST("/create", logistic.Create)
				logistics.PUT("/update", logistic.Update)
			}
			// 运费模板
			freights := setup.Group("/freight")
			{
				freights.GET("/search", freight.Search)
				freights.GET("/id", freight.Get)
				freights.DELETE("/delete/:id", freight.Delete)
				freights.POST("/create", freight.Create)
				freights.PUT("/update", freight.Update)

				// 运费模板--城市
				freight_citys := freights.Group("/city")
				{
					freight_citys.GET("/search", freight_city.Search)
					freight_citys.GET("/id", freight_city.Get)
					freight_citys.DELETE("/delete/:id", freight_city.Delete)
					freight_citys.POST("/create", freight_city.Create)
					freight_citys.POST("/createAll", freight_city.CreateAll)
					freight_citys.PUT("/update", freight_city.Update)
				}

			}
		}

		// 小程序
		wxs := v.Group("/wx")
		{
			wxs.POST("/login", wx.Login)
			wxs.GET("/info", wx.Info)
			wxs.GET("/phone", wx.Phone)
			wxs.POST("/pay", wx.Pay)
			wxs.POST("/notify/pay", wx.PayNotify)
			wxs.POST("/notify/refund/:secret", wx.RefundNotify)
			wxs.GET("/access_token", wx.GetAccessToken)
			// 小程序码
			wxs.GET("/qrcode", wx.GetQRCode)
			wxs.GET("/qrcode/key", wx.GetByKey)
			// 普通二维码
			wxs.GET("/prcode", wx.PQrcode)

			// 退款
			wxs.POST("/refund", wx.Refund)
		}

		alis := v.Group("/ali")
		{
			smss := alis.Group("/sms")
			{
				smss.POST("/send", sms.Send)
				smss.POST("/check", sms.Check)
			}
		}

		// 模块开启
		modules := v.Group("/module")
		{
			modules.GET("/search", module.GetBySearch)
			modules.GET("/id", module.GetByID)
			modules.DELETE("/delete/:id", module.Delete)
			modules.POST("/create", module.Create)
			modules.PUT("/update", module.Update)
		}
		// 初始化模块路由
		ModuleRouter()

		// 日志记录
		records := v.Group("/record")
		{
			records.GET("/search", record.Search)
		}

		// pca
		pcas := v.Group("/pca")
		{
			pcas.GET("/ws", pca.Ws)
		}
	}
}
