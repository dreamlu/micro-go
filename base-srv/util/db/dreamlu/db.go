package dreamlu

import (
	"demo/base-srv/controllers/wx"
	"demo/base-srv/models/admin"
	"demo/base-srv/models/admin/applet"
	"demo/base-srv/models/admin/setup/freight"
	"demo/base-srv/models/admin/setup/freight/freight_city"
	"demo/base-srv/models/admin/setup/logistic"
	"demo/base-srv/models/carousel"
	"demo/base-srv/models/client"
	"demo/base-srv/models/client/address"
	"demo/base-srv/models/client/cart"
	"demo/base-srv/models/goods"
	"demo/base-srv/models/goods/category"
	"demo/base-srv/models/goods/norm"
	"demo/base-srv/models/goods/notice"
	"demo/base-srv/models/goods/regroup"
	"demo/base-srv/models/goods/warehouse"
	"demo/base-srv/models/module"
	"demo/base-srv/models/module/day_sign"
	"demo/base-srv/models/module/day_sign/img"
	"demo/base-srv/models/module/day_sign/text"
	"demo/base-srv/models/module/live_bro"
	"demo/base-srv/models/order"
	"demo/base-srv/models/order/comment"
	"demo/base-srv/models/order/order_refund"
	"demo/base-srv/models/order/order_refund/info"
	"demo/base-srv/models/record"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/util"
)

func InitDB() {

	gt.NewCrud().DB().AutoMigrate(
		&carousel.Carousel{},
		&info.OrderRInfo{},
		&order_refund.OrderRefund{},
	)
	gt.NewCrud().DB().AutoMigrate(
		&client.Client{},
		&admin.Admin{},
		&order.Order{},
		&applet.Applet{},
		&day_sign.DaySign{},
		&img.DaySignImg{},
		&text.DaySignText{},
		&cart.Cart{},
		&order.OrderGoods{},
		&wx.QrCode{},
	)
	gt.NewCrud().DB().AutoMigrate(
		&category.GsCategory{},
		&comment.OrderComment{},
		&logistic.Logistic{},
		&freight.Freight{},
		&freight_city.FreightCity{},
		&norm.GsNorm{},
		&live_bro.LiveBroData{},
		&record.Record{},
		&regroup.GsRegroup{},
		&notice.GsNotice{},
		&warehouse.GsWarehouse{},
		&module.Module{},
	)
	gt.NewCrud().DB().AutoMigrate(
		&address.ClientAddress{},
		&goods.Goods{},
	)
	initSQL()
}

func initSQL() {
	alterTable()
	// 插入admin账号
	role := int8(1)
	var ad = admin.Admin{
		Name:     "admin",
		Password: util.AesEn("123456"),
		Role:     &role,
	}
	gt.NewCrud(gt.Data(&ad), gt.Model(admin.Admin{})).Create()
}

func alterTable() {
	//gt.NewCrud().Select("alter table `order` modify order_num varchar(30) null;").Exec()
}
