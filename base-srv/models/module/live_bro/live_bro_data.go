package live_bro

import (
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
	"strconv"
)

// 小程序直播模型
type LiveBroData struct {
	models2.AdminCom
	Data json.CJSON `json:"data" gorm:"type:json"` // 直播列表数据
	//ReplyData json.CJSON // 回放数据
}

var crud = gt.NewCrud(
	gt.Model(LiveBroData{}),
)

// get data, by id
func (c *LiveBroData) Get(admin_id uint64) (*LiveBroData, error) {

	var (
		params = cmap.CMap{}
		data   LiveBroData // not use *OrderComment
	)
	params.Add("admin_id", strconv.FormatUint(admin_id, 10))
	crud.Params(gt.Data(&data))
	if err := crud.GetByData(params).Error(); err != nil {
		return nil, err
	}
	return &data, nil
}

// update data
func (c *LiveBroData) Update(data *LiveBroData) (*LiveBroData, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

// create data
func (c *LiveBroData) Create(data *LiveBroData) (*LiveBroData, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}

	return data, nil
}
