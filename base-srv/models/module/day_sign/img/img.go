package img

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	models2 "micro-go/commons/models"
)

// 日签--图片模型
type DaySignImg struct {
	models2.ModelCom
	DaySignID uint   `json:"day_sign_id" gorm:"type:bigint(20);INDEX:日签图片查询索引"` // 日签id
	ImgUrl    string `json:"img_url"`                                           // 图片url
}

var crud = gt.NewCrud(
	gt.Model(DaySignImg{}),
)

// get data, by id
func (c *DaySignImg) GetByID(id string) interface{} {

	var data DaySignImg // not use *DaySignImg
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *DaySignImg) GetBySearch(params map[string][]string) interface{} {
	var datas []*DaySignImg
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *DaySignImg) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *DaySignImg) Update(data *DaySignImg) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *DaySignImg) Create(data *DaySignImg) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
