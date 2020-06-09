package day_sign

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	models2 "micro-go/commons/models"
)

// 日签模型
type DaySign struct {
	models2.AdminCom
	ImgUrl  string `json:"img_url"` // 封面
	Content string `json:"content"` // 内容
}

var crud = gt.NewCrud(
	gt.Model(DaySign{}),
)

// get data, by id
func (c *DaySign) GetByID(id string) interface{} {

	var data DaySign // not use *DaySign
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *DaySign) GetBySearch(params map[string][]string) interface{} {
	var datas []*DaySign
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *DaySign) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *DaySign) Update(data *DaySign) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *DaySign) Create(data *DaySign) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
