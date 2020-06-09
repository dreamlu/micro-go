package order_refund

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
	"micro-go/base-srv/models/order"
	models2 "micro-go/commons/models"
)

// 订单退款模型
type OrderRefund struct {
	models2.AdminCom
	ClientID     uint64     `json:"client_id" gorm:"type:bigint(20)"`              // 客户id
	OrderID      uint64     `json:"order_id" gorm:"type:bigint(20)"`               // 订单id
	OrderGoodsID uint64     `json:"order_goods_id" gorm:"type:bigint(20)"`         // 订单商品id
	GoodsStatus  *int8      `json:"goods_status" gorm:"type:tinyint(2);DEFAULT:0"` // 货物状态
	Reason       *int8      `json:"reason" gorm:"type:tinyint(2);DEFAULT:0"`       // 退款原因
	Remark       string     `json:"remark"`                                        // 退款备注
	Img          json.CJSON `json:"img" gorm:"type:json"`                          // 凭证图片,json格式
	Reply        string     `json:"reply"`                                         // 商家回复

	// update
	// 4退款完成,5申请退款中(默认),6拒绝退款
	Status int8 `json:"status" gorm:"type:tinyint(2);DEFAULT:5"`
	//OrderStatus int8    `json:"order_status" gorm:"type:tinyint(2);DEFAULT:5"` // 订单本身状态,记录用
	GoodsID     uint64  `json:"goods_id" gorm:"type:bigint(20)"`                                     // 商品id
	Name        string  `json:"name" gorm:"type:varchar(50)"`                                        // 名称
	Num         int     `json:"num"`                                                                 // 数量
	ImgUrl      string  `json:"img_url"`                                                             // 封面
	Money       float64 `json:"money" gorm:"type:double(10,2)"`                                      // 申请退款金额
	OutRefundNo string  `json:"out_refund_no" gorm:"type:bigint(20)"`                                // 唯一退款单号
	OutTradeNo  string  `json:"out_trade_no" gorm:"type:bigint(20)"`                                 // 支付单号
	OrderNum    string  `json:"order_num" gorm:"type:varchar(50)" gt:"field:order_refund.order_num"` // 订单编号
	//Info        json.CJSON `json:"info" gorm:"type:json"`                                               // 已收到货,填写的商家地址等信息
}

type OrderRefundGs struct {
	OrderRefund
	OrderStatus int8                   `json:"order_status" ` // 订单本身状态
	ClientName  string                 `json:"client_name"`
	OrderGoods  []*order.OrderGoodsPar `json:"order_goods" gt:"sub_sql"`

	// 已收到货, 退货信息
	OrderRInfoStatus byte `json:"order_r_info_status" gt:"field:order_r_info.status"` // 货物状态
	//OrderRInfoInfo   string `json:"order_r_info_info"`
}

type Num struct {
	Num int `json:"num"`
}

var crud = gt.NewCrud(
	gt.Model(OrderRefund{}),
)

// id
func (c *OrderRefund) FindByID(id uint64) error {

	crud.Params(gt.Data(&c))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return err
	}
	return nil
}

// get data, by id
func (c *OrderRefund) GetByID(id interface{}) (data OrderRefund, err error) {

	crud.Params(gt.Data(&data))
	if err = crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return
	}
	return
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *OrderRefund) GetBySearch(params cmap.CMap) interface{} {

	keySQL := ""
	order_r_info_statusSQL := ""
	if key := params.Get("key"); key != "" {
		keySQL = "order_r_info.info like '%" + key + "%'"
		params.Del("key")
	}
	if order_r_info_status := params.Get("order_r_info_status"); order_r_info_status != "" {
		order_r_info_statusSQL = "order_r_info.status = " + order_r_info_status
		params.Del("order_r_info_status")
	}
	var datas []*OrderRefundGs
	cd := gt.NewCrud(
		gt.Model(OrderRefundGs{}),
		gt.InnerTable([]string{"order_refund", "client", "order_refund", "order"}),
		gt.LeftTable([]string{"order_refund:id", "order_r_info:order_refund_id"}),
		gt.Data(&datas),
		gt.SubWhereSQL(keySQL, order_r_info_statusSQL),
	)
	cd = cd.GetMoreBySearch(params)
	if cd.Error() != nil {
		return result.CError(cd.Error())
	}
	// 查询关联商品订单
	var or order.Order
	for _, v := range datas {
		ords, err := or.GetGoods(v.OrderID)
		if err != nil {
			return result.CError(err)
		}
		v.OrderGoods = ords
	}

	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *OrderRefund) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
// 指针相互引用问题
func (c OrderRefund) Update(data *OrderRefund) (err error) {

	cd := gt.NewCrud().Begin()
	cd.Params(gt.Data(data))
	if err = cd.Update().Error(); err != nil {
		cd.Rollback()
		return
	}

	if data.Status != 0 {
		//fmt.Println("orderRefund之前数据: ", data)
		err = c.FindByID(data.ID)
		if err != nil {
			cd.Rollback()
			return
		}
		//fmt.Println("orderRefund之后数据: ", c)
		if err = cd.Select("update `order_goods` set status = ? where id = ?", data.Status, c.OrderGoodsID).Exec().Error(); err != nil {
			gt.Logger().Error(err)
			cd.Rollback()
			return
		}
	}

	// 申请通过
	if data.Status == 4 {
		fmt.Println("[通过开始判断退款是否全部通过]")
		// 判断同一订单的商品全部都退款审核通过
		var (
			ids []uint64
			or  order.Order
		)
		ogs, err := or.GetGoodsID(c.OrderID)
		if err != nil {
			cd.DB().Rollback()
			return err
		}
		for _, v := range ogs {
			ids = append(ids, v.ID)
		}

		gt.Logger().Info("[退款申请通过的ids]:", ids)
		var num Num
		cd2 := gt.NewCrud(gt.Data(&num)).Select("select count(*) num from `order_refund` where order_goods_id in (?) and status = 4", ids).Single()
		if cd2.Error() != nil {
			gt.Logger().Error("[订单退款申请判断是否全部通过问题]:", cd2.Error())
			cd.Rollback()
			return cd2.Error()
		}
		gt.Logger().Info("[退款申请通过的行数]:", num.Num)
		// 都退款通过
		if num.Num == len(ids)-1 { // 事务查询影响 - 1
			gt.Logger().Info("[修改订单整体退款状态,订单id]:", c.OrderID)
			db := cd.DB().Exec("update `order` set status = 4 where id = ?", c.OrderID)
			if db.Error != nil {
				gt.Logger().Error(db.Error)
				cd.Rollback()
				return db.Error
			}
			go func() {
				// 删除对应分销
				gt.NewCrud().Select("delete from dis.income where order_id = ?", c.OrderID).Exec()
			}()
		}
	}

	if err = cd.Commit().Error(); err != nil {
		return
	}
	return
}

// create data
func (c *OrderRefund) Create(data *OrderRefund) interface{} {

	cd := crud.Begin()

	newID, _ := id.NewID(1)
	data.OutRefundNo = newID.String()

	cd.Params(gt.Data(data))
	if err := cd.Create().Error(); err != nil {
		cd.DB().Rollback()
		return result.CError(err)
	}

	// 申请退款, 状态5
	if err := cd.Select("update `order_goods` set status = 5 where id = ?", data.OrderGoodsID).Exec().Error(); err != nil {
		cd.DB().Rollback()
		return result.CError(err)
	}

	if err := cd.Commit().Error(); err != nil {
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
