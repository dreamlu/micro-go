package order

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/time"
	"micro-go/base-srv/models/record"
	"micro-go/base-srv/util/cons"
	time2 "time"
)

// 订单更新后的一些操作
func otherUp(cd gt.Crud, data *Order) error {

	// 订单之前和之后的状态
	var befData, afData Order
	gt.NewCrud().DB().First(&befData, data.ID)
	cd.DB().First(&afData, data.ID)

	if befData.ID == 0 || afData.ID == 0 {
		return result.TextError("[订单更新之前或之后状态数据不存在]")
	}

	// 一.购买
	if *befData.Status == 0 && *afData.Status == 1 {
		// 1.商品购买量增加
		go addGoodsBuyNum(cd, &afData)
		// 2.创建分销计算与分销记录
		//go Distr(&afData)
	}

	// 二.对应收入(分销)记录修改
	if *befData.Status >= 1 && *afData.Status <= 4 {
		//go disIncomeStatus(&afData)
	}

	// 三.订单日志
	go orderLog(befData, afData)

	return nil
}

// 该订单下的所有商品
// 商品购买量增加
func addGoodsBuyNum(cd gt.Crud, data *Order) {

	// 该笔订单下所有购买的商品
	gs, err := data.GetGoods(data.ID)
	if err != nil {
		gt.Logger().Error(err.Error())
		return
	}

	// 增加对应规格已售数量
	crud := gt.NewCrud()
	for _, v := range gs {
		if v.GsNormID == 0 {
			crud.Select("update `goods` set sell_num = sell_num + ? where id = ?", v.Num, v.GoodsID).Exec()
		} else {
			crud.Select("update `gs_norm` set sell_num = sell_num + ? where id = ?", v.Num, v.GsNormID).Exec()
		}
	}
}

// 操作记录
func orderLog(befData Order, afData Order) {
	// 购买
	befS := Status(*befData.Status)
	afS := Status(*afData.Status)
	if befS != afS {
		record.Log(cons.Order, afData.AdminID, "客户["+afData.ClientName+"]订单状态由["+befS+"]变成["+afS+"];时间:"+time.CTime(time2.Now()).String()+";支付编号:"+afData.OutTradeNo, afData)
	}
}

func Status(status int8) (res string) {
	switch status {
	case 0:
		res = "待付款"
	case 1:
		res = "待发货"
	case 2:
		res = "待收货"
	case 3:
		res = "已完成"
	case 4:
		res = "交易关闭" //退款完成"
	case 5:
		res = "退款申请中"
	case 6:
		res = "拒绝退款"
	case 7:
		res = "待评价"
	case 8:
		res = "订单删除"
	default:
		res = "未知状态"
	}
	return
}
