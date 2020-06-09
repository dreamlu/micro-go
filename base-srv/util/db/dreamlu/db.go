package dreamlu

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/util"
	"micro-go/base-srv/controllers/wx"
	"micro-go/base-srv/models/admin"
	"micro-go/base-srv/models/admin/applet"
	"micro-go/base-srv/models/admin/setup/freight"
	"micro-go/base-srv/models/admin/setup/freight/freight_city"
	"micro-go/base-srv/models/admin/setup/logistic"
	"micro-go/base-srv/models/carousel"
	"micro-go/base-srv/models/client"
	"micro-go/base-srv/models/client/address"
	"micro-go/base-srv/models/client/cart"
	"micro-go/base-srv/models/goods"
	"micro-go/base-srv/models/goods/category"
	"micro-go/base-srv/models/goods/norm"
	"micro-go/base-srv/models/goods/notice"
	"micro-go/base-srv/models/goods/regroup"
	"micro-go/base-srv/models/goods/warehouse"
	"micro-go/base-srv/models/module"
	"micro-go/base-srv/models/module/day_sign"
	"micro-go/base-srv/models/module/day_sign/img"
	"micro-go/base-srv/models/module/day_sign/text"
	"micro-go/base-srv/models/module/live_bro"
	"micro-go/base-srv/models/order"
	"micro-go/base-srv/models/order/comment"
	"micro-go/base-srv/models/order/order_refund"
	"micro-go/base-srv/models/order/order_refund/info"
	"micro-go/base-srv/models/record"
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
