package carousel

import (
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
)

// 轮播图模型
type Carousel struct {
	models2.AdminCom
	Name   string `gorm:"type:varchar(30)" json:"name" valid:"required,len=2-20"` // 名称
	ImgUrl string `json:"img_url"`                                                // 图片url
	Url    string `json:"url"`                                                    // 跳转地址
	Type   byte   `json:"type" gorm:"type:tinyint(1);DEFAULT:0"`                  // 0商品详情(默认),1直播
}

var crud = gt.NewCrud(
	gt.Model(Carousel{}),
)

// get data, by id
func (c *Carousel) GetByID(id string) interface{} {

	var data Carousel // not use *Carousel
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Carousel) GetBySearch(params map[string][]string) interface{} {
	var datas []*Carousel
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *Carousel) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *Carousel) Update(data *Carousel) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *Carousel) Create(data *Carousel) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
