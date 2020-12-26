package module

import (
	"demo/base-srv/models/admin/applet"
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"strconv"
)

// 商家开启模块 模型
type Module struct {
	models2.AdminCom
	M int8  `json:"m" gorm:"tinyint(3)"` // 0日签,1直播,2优惠券,3短视频
	T *int8 `json:"t"`                   // 类型
	//IsOpen int8 `json:"is_open" gorm:"tinyint(2);DEFAULT:0"` // 0关闭(默认),1开启
}

var crud = gt.NewCrud(
	gt.Model(Module{}),
)

// get data, by id
func (c *Module) GetByID(id string) interface{} {

	var data Module // not use *Module
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Module) GetBySearch(params cmap.CMap) interface{} {

	var wx applet.Applet
	wx.Appid = params.Get("appid")
	if wx.Appid != "" {
		params.Del("appid")
		if err := wx.GetByAppid(wx.Appid); err != nil {
			return result.CError(err)
		}
		params.Add("admin_id", strconv.FormatUint(wx.AdminID, 10))
	}

	var datas []*Module
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager()).Add("admin_id", wx.AdminID).Add("logo", wx.Logo)
}

// delete data, by id
func (c *Module) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *Module) Update(data *Module) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *Module) Create(data *Module) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate).Add("id", data.ID)
}
