package order

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
	"github.com/dreamlu/gt/tool/type/time"
	"micro-go/base-srv/models/goods"
	"micro-go/base-srv/models/record"
	"micro-go/base-srv/util/cons"
	models2 "micro-go/commons/models"
	"strconv"
	time2 "time"
)

// Order
type Order struct {
	models2.AdminCom
	ClientID uint64 `json:"client_id" gorm:"type:bigint(20);INDEX:查询索引client_id"` // 客户id
	// 0待付款,1待发货,2待收货,3已完成(待评价/确认收货),
	// 4退款完成(交易取消),5申请退款中(废弃),6拒绝退款(废弃),
	// 7已评价,8已取消
	Status         *int8      `json:"status" gorm:"type:tinyint(2);DEFAULT:0"`
	Money          float64    `json:"money" gorm:"type:double(10,2)"`                   // 付款金额
	OutTradeNo     string     `json:"out_trade_no" gorm:"type:varchar(50)"`             // 商户订单号(支付订单号, 退款等用)
	ClientName     string     `json:"client_name" gorm:"type:varchar(50)"`              // 姓名
	Phone          string     `json:"phone" gorm:"type:varchar(20)"`                    // 手机号
	Address        string     `json:"address"`                                          // 地址
	FreightMoney   float64    `json:"freight_money" gorm:"type:double(10,2);DEFAULT:0"` // 运费,默认0
	Remark         string     `json:"remark"`                                           // 备注
	Paytime        time.CTime `json:"paytime" gorm:"type:datetime"`                     // 支付时间,后端创建
	Shiptime       time.CTime `json:"shiptime" gorm:"type:datetime"`                    // 确认收货时间,前端创建
	LogisticInfo   json.CJSON `json:"logistic_info" gorm:"type:json"`                   // 物流信息,json
	ClientCouponID uint64     `json:"client_coupon_id"`                                 // 客户优惠券id
	CouponMoney    float64    `json:"coupon_money" gorm:"type:double(10,2)"`            // 优惠券减免金额
	// update
	OrderNum string `json:"order_num" gorm:"type:varchar(30)"` // 订单号(同一个货仓的同一笔支付)
}

type OrderGS struct {
	Order
	OrderGoods []*OrderGoodsPar `json:"order_goods"`
}

type OrderGSID struct {
	Order
	OrderGoods []*OrderGoods `json:"order_goods"`
}

var crud = gt.NewCrud(
	gt.Model(Order{}),
)

func (c *Order) GetByOrderNum(orderNum string) (data Order, err error) {

	gt.Logger().Info("[根据outTradeNo批量查询订单]")
	var params = cmap.CMap{}
	params.Add("`order_num`", orderNum)
	crud2 := gt.NewCrud()
	crud2.Params(gt.Model(Order{}), gt.Table("order"), gt.Data(&data))
	if err = crud2.GetByData(params).Error(); err != nil {
		gt.Logger().Error(err.Error())
		return
	}
	return
}

func (c *Order) GetByOutTradeNo(outTradeNo string) (datas []*Order, err error) {

	gt.Logger().Info("[根据outTradeNo批量查询订单]")
	var params = cmap.CMap{}
	params.Add("`out_trade_no`", outTradeNo)
	crud := gt.NewCrud()
	crud.Params(gt.Model(Order{}), gt.Table("order"), gt.Data(&datas))
	if err = crud.GetByData(params).Error(); err != nil {
		gt.Logger().Error(err.Error())
		return
	}
	return
}

// id
func (c *Order) FindByID(id interface{}) error {

	crud.Params(gt.Data(&c))
	if err := crud.GetByID(id).Error(); err != nil {
		return err
	}
	return nil
}

// get data, by id
func (c *Order) GetByID(id interface{}) (datags OrderGSID, err error) {

	var (
		data Order
		//datags OrderGS // not use *Order
	)
	crud.Params(gt.Data(&data))
	err = crud.GetByID(id).Error()
	if err != nil {
		return
	}
	datags.Order = data
	datags.OrderGoods, err = c.GetGoodsID(id)
	if err != nil {
		return
	}
	return
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Order) GetBySearch(params cmap.CMap) interface{} {

	statusWhereSQL := ""
	status := params.Get("status")
	switch {
	case status == "refund":
		params.Del("status")
		statusWhereSQL = "status in (4,5,6)"
	case status == "comment":
		params.Del("status")
		statusWhereSQL = "status in (3,7)"
	case status != "":
		params.Del("status")
		statusWhereSQL = "status in (" + status + ")"
	}

	// 查询基础订单
	var datas []*OrderGSID
	crud2 := gt.NewCrud(
		gt.Table("order"),
		gt.Model(Order{}),
		gt.Data(&datas),
		gt.SubWhereSQL(statusWhereSQL),
	)
	cd := crud2.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	// 查询关联商品订单
	for _, v := range datas {
		ords, err := c.GetGoodsID(v.ID)
		if err != nil {
			return result.CError(err)
		}
		v.OrderGoods = ords
	}

	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *Order) Delete(id string) interface{} {

	var befData Order
	gt.NewCrud().DB().First(&befData, id)

	if err := crud.Delete(id).Error(); err != nil {
		return result.GetError(err)
	}
	// 删除订单关联的其他表信息
	go func() {
		crud.Select("delete from `order_goods` where order_id = ?", id).Exec()
		// 操作记录
		go record.Log(cons.Order, befData.AdminID, "客户["+befData.ClientName+"]订单被取消, 订单支付编号:"+befData.OutTradeNo, befData)
	}()
	return result.MapDelete
}

// 批量通知修改
// 和update内容相同,只是批量罢了
func (c *Order) UpdateNotify(data *Order) (err error) {
	datas, err := c.GetByOutTradeNo(data.OutTradeNo)
	if err != nil {
		return
	}
	cd := crud.Begin()
	for _, v := range datas {
		data.ID = v.ID

		cd.Params(gt.Data(data))
		if err = cd.Update().Error(); err != nil {
			cd.Rollback()
			return
		}

		err = otherUp(cd, data)
		if err != nil {
			cd.Rollback()
			return
		}
	}
	if err = cd.Commit().Error(); err != nil {
		return
	}
	return nil
}

//// update data
//func (c *Order) UpdateMore(datas []*Order) error {
//
//	for _, v := range datas {
//		_, err := c.Update(v)
//		if err != nil {
//			gt.Logger().Error("[订单批量修改错误]", err.Error())
//		}
//	}
//	return nil
//}

// update data
func (c *Order) Update(data *Order) (*Order, error) {

	cd := crud.Begin()
	cd.Params(gt.Data(data))
	if err := cd.Update().Error(); err != nil {
		cd.Rollback()
		return nil, err
	}

	err := otherUp(cd, data)
	if err != nil {
		cd.Rollback()
		return nil, err
	}

	if err := cd.Commit().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

// 订单分割
// 与实际支付价格
func splitOrder(data *OrderGS) (datas []*OrderGS) {

	//if len(data.OrderGoods) == 1 {
	//	datas = append(datas, data)
	//	return
	//}
	// 同仓库商品拆分
	for _, v := range data.OrderGoods {
		if b := exits(v, datas); !b {
			// 创建发货仓
			var dataT OrderGS
			dataT.Order = data.Order

			// 优化
			//if len(data.OrderGoods) > 2 {
			dataT.Money = 0
			dataT.CouponMoney = 0
			//dataT.OrderGoods = data.OrderGoods // 单个商品, 无需进行订单分割
			addGoods(v.WarehouseID, *data, &dataT)
			//} else {

			//}
			datas = append(datas, &dataT)
		}
	}
	return
}

// 1.将同一发货仓商品加入同一订单中
// 2.邮费计算
// 3.同时进行实际金钱的计算分割
func addGoods(wid uint64, data OrderGS, newOr *OrderGS) {
	var allPrice float64
	// 加入订单
	for _, v := range data.OrderGoods {
		if v.WarehouseID == wid {
			newOr.OrderGoods = append(newOr.OrderGoods, v)
		}
		allPrice += v.Price * float64(v.Num)
	}
	newOr.FreightMoney = maxFreightMoney(newOr)
	// 计算时不需要邮费 实际金额 - 总邮费
	data.Money -= data.FreightMoney
	// 后续有单个商品优惠券时, 可以商品实际价格 = 商品实际价格 - 商品优惠券(添加字段记录)
	for _, v := range newOr.OrderGoods {
		pro := v.Price * float64(v.Num) / allPrice
		v.Money = pro * data.Money
		newOr.Money += v.Money
		newOr.CouponMoney += data.CouponMoney * pro
	}
	// 实际结果加上邮费
	newOr.Money += newOr.FreightMoney
	newID, _ := id.NewID(1)
	nano := strconv.FormatInt(time2.Now().UnixNano(), 10)
	nano = string([]byte(nano)[len(nano)-5:])
	newOr.OrderNum = fmt.Sprint(newID.String(), nano)
}

// 获取同一订单商品中最大邮费
func maxFreightMoney(newOr *OrderGS) (freightMoney float64) {

	freightMoney = newOr.OrderGoods[0].FreightMoney
	for _, v := range newOr.OrderGoods {
		if v.FreightMoney > freightMoney {
			freightMoney = v.FreightMoney
		}
	}
	return
}

// 判断统一发货仓是否已存在
func exits(data *OrderGoodsPar, datas []*OrderGS) (b bool) {
	for _, v2 := range datas {
		// 一个订单至少存在一个商品
		if v2.OrderGoods[0].WarehouseID == data.WarehouseID {
			return true
		}
	}
	return
}

// create data
func (c *Order) Create(data *OrderGS) (err error) {

	// 验证,重复下单等
	err = check(data)
	if err != nil {
		return err
	}

	// 根据发货仓进行订单分割
	datas := splitOrder(data)

	cd := crud.Begin()
	for _, v := range datas {
		err = create(cd, v)
		if err != nil {
			cd.Rollback()
			return err
		}
	}

	// 优惠券使用后状态
	// 简单点
	err = cd.Select("update coupon.client_coupon set `status` = 1 where id = ?", data.ClientCouponID).Exec().Error()
	if err != nil {
		cd.Rollback()
		return err
	}

	if err = cd.Commit().Error(); err != nil {
		return err
	}
	return
}

// 单个订单创建
func create(cd gt.Crud, data *OrderGS) (err error) {

	or, err := createCom(cd, &data.Order)
	if err != nil {
		return
	}
	for k := range data.OrderGoods {
		data.OrderGoods[k].OrderID = or.ID
	}
	if len(data.OrderGoods) > 0 {
		cd.Params(
			gt.Model(OrderGoodsPar{}),
			gt.Table("order_goods"),
			gt.Data(data.OrderGoods),
		)
		if err = cd.CreateMore().Error(); err != nil {
			return
		}
	}

	return
}

// 基本订单创建
func createCom(cd gt.Crud, data *Order) (*Order, error) {
	// 基本订单创建
	cd.Params(
		gt.Model(Order{}),
		gt.Table("order"),
		gt.Data(data),
	)
	if err := cd.Create().Error(); err != nil {
		cd.DB().Rollback()
		gt.Logger().Error(err.Error())
		return nil, err
	}
	return data, nil
}

type Num struct {
	Num int `json:"num"`
}

func check(data *OrderGS) (err error) {

	// 判断重复下单:redis
	key := strconv.FormatUint(data.ClientID, 10) + strconv.FormatUint(data.OrderGoods[0].GoodsID, 10)
	ce := cache.NewCache()
	ca, _ := ce.Get(key)
	//if err != nil &&  {
	//	return err
	//}
	if ca.Data == nil {
		ca.Data = 1
		ca.Time = 5 * cache.CacheSecond
		_ = ce.Set(key, ca)
		//return nil
	} else {
		return result.TextError("重复下单")
	}

	// 商品验证
	var gs goods.Goods
	for _, v := range data.OrderGoods {
		// 1.进行时间判断
		// 查找原商品时间
		g, err := gs.Get(v.GoodsID)
		if err != nil {
			return err
		}
		//if time2.Now().After(time2.Time(g.EndTime)) {
		//	return result.TextError("[商品已过期]")
		//}

		// 2.库存判断, 单个规格剩余数量 = 商品总库存 - 已卖数量 - 待付款数量
		var num Num
		cd := gt.NewCrud(gt.Data(&num))
		if v.GsNormID == 0 { // 单规格
			cd.Select("select sum(bb.num) as num from `order` aa inner join `order_goods` bb on aa.id = bb.order_id where aa.status = 0 and goods_id = ?", v.GoodsID).
				Single()
		} else { // 多规格
			cd.Select("select sum(bb.num) as num from `order` aa inner join `order_goods` bb on aa.id = bb.order_id where aa.status = 0 and goods_id = ? and bb.gs_norm_id = ?", v.GoodsID, v.GsNormID).
				Single()
		}
		// 总销售量
		for _, v := range g.Norm {
			g.SellNum += v.SellNum
			g.Inventory += v.Inventory
		}
		leftNum := g.Inventory - g.SellNum - num.Num
		if leftNum <= 0 {
			return result.TextError("存在商品库存不足")
		}
	}

	return
}

// map
//key -> value

//sql where key =?
//sql where key = value
