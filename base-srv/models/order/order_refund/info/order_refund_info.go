package info

import (
	"demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
)

// 订单申请退款--已收到货 模型
type OrderRInfo struct {
	models.ModelCom
	OrderRefundID uint64 `json:"order_refund_id" gorm:"UNIQUE_INDEX:已存在退款发货信息"` // 订单退款id
	// 0nothing,1商家同意,填地址(退货申请,默认), 2客户填写信息(退货中), 3(退货完成)
	Status  byte       `json:"status" gorm:"type:tinyint(2);DEFAULT:1"`
	Address string     `json:"address"`               // 地址
	Info    json.CJSON `json:"info" gorm:"type:json"` // 客户填写信息
}

var crud = gt.NewCrud(
	gt.Model(OrderRInfo{}),
)

// get data, by id
func (c *OrderRInfo) Get(params cmap.CMap) (data OrderRInfo, err error) {
	crud.Params(gt.Data(&data))
	if err = crud.GetByData(params).Error(); err != nil {
		return
	}
	return
}

// get data, limit and search
// OrderRInfoPage 1, everyPage 10 default
func (c *OrderRInfo) Search(params cmap.CMap) (datas []*OrderRInfo, pager result.Pager, err error) {

	crud.Params(
		gt.Data(&datas),
	)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// delete data, by id
func (c *OrderRInfo) Delete(id interface{}) error {

	return crud.Delete(id).Error()
}

// update data
func (c *OrderRInfo) Update(data *OrderRInfo) (*OrderRInfo, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *OrderRInfo) Create(data *OrderRInfo) (*OrderRInfo, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}

	return data, nil
}
