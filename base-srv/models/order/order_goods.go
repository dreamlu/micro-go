package order

import (
	models2 "demo/commons/models"
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/type/json"
)

// 订单--商品 模型
type OrderGoods struct {
	models2.IDCom
	OrderGoodsPar
}

// 商品列表
type OrderGoodsPar struct {
	OrderID  uint64     `json:"order_id" gorm:"type:bigint(20);INDEX:订单商品列表查询order_id索引"` // 订单id
	GoodsID  uint64     `json:"goods_id" gorm:"type:bigint(20)"`                          // 商品id
	Name     string     `json:"name" gorm:"type:varchar(50)"`                             // 名称
	Num      int        `json:"num"`                                                      // 数量
	ImgUrl   string     `json:"img_url"`                                                  // 封面
	Price    float64    `json:"price" gorm:"type:double(10,2)"`                           // 商品价格
	Money    float64    `json:"money" gorm:"type:double(10,2)"`                           // 实际支付价格,后端计算
	GsNorm   json.CJSON `json:"gs_norm" gorm:"type:json"`                                 // 规格记录, json
	GsNormID uint64     `json:"gs_norm_id" gorm:"type:bigint(20);INDEX:查询索引client_id"`    // 规格id

	// 社群特有
	HeadPrice   float64 `json:"head_price" gorm:"type:double(10,2)"`    // 团长价
	SpHeadPrice float64 `json:"sp_head_price" gorm:"type:double(10,2)"` // 超级团长价

	// update
	WarehouseID   uint64  `json:"warehouse_id" gorm:"type:bigint(20)"`              // 货仓id
	WarehouseName string  `json:"warehouse_name" gorm:"type:varchar(50)"`           // 货仓名称
	FreightMoney  float64 `json:"freight_money" gorm:"type:double(10,2);DEFAULT:0"` // 运费,默认0

	// 0nothing(默认),4退款完成,5申请退款中,6拒绝退款
	Status int8 `json:"status" gorm:"type:tinyint(2);DEFAULT:0"`
}

func (c *Order) GetGoods(orderID interface{}) ([]*OrderGoodsPar, error) {
	var datas []*OrderGoodsPar
	cd := crud.Params(gt.Data(&datas)).
		Select(fmt.Sprintf("select %s from order_goods where order_id = ?", gt.GetColSQL(OrderGoodsPar{})), orderID).Single()
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return nil, cd.Error()
	}
	return datas, nil
}

func (c *Order) GetGoodsID(orderID interface{}) ([]*OrderGoods, error) {
	var datas []*OrderGoods
	cd := crud.Params(gt.Data(&datas)).
		Select(fmt.Sprintf("select %s from order_goods where order_id = ?", gt.GetColSQL(OrderGoods{})), orderID).Single()
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return nil, cd.Error()
	}
	return datas, nil
}

func (c *OrderGoods) FindByID(id interface{}) (err error) {

	crud.Params(gt.Data(&c))
	err = crud.GetByID(id).Error()
	if err != nil {
		return
	}
	return
}

//
//func (c *OrderGoods) GetByOrderID(orderID interface{}) (datas []*OrderGoods,err error) {
//	cd := gt.NewCrud(gt.Model(OrderGoods{}), gt.Data(&datas))
//	if err = cd.Select(fmt.Sprintf("select %s from `order_goods` where order_id")).Single().Error(); err != nil {
//		return
//	}
//	return
//}
