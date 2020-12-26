package text

import (
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
)

// 日签--文字模型
type DaySignText struct {
	models2.ModelCom
	DaySignID uint   `json:"day_sign_id" gorm:"type:bigint(20);INDEX:日签图片查询索引"` // 日签id
	Content   string `json:"content"`                                           // 文字
}

var crud = gt.NewCrud(
	gt.Model(DaySignText{}),
)

// get data, by id
func (c *DaySignText) GetByID(id string) interface{} {

	var data DaySignText // not use *DaySignText
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *DaySignText) GetBySearch(params map[string][]string) interface{} {
	var datas []*DaySignText
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *DaySignText) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *DaySignText) Update(data *DaySignText) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *DaySignText) Create(data *DaySignText) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
