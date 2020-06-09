package admin

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/util"
	models2 "micro-go/commons/models"
)

// 多账号管理 模型
type Admin struct {
	models2.ModelCom
	Name     string `gorm:"type:varchar(30);UNIQUE_INDEX:账号已存在" json:"name"` // 名称
	Password string `json:"password" gorm:"type:varchar(100)"`               // 密码
	Role     *int8  `json:"role" gorm:"type:tinyint(2);DEFAULT:0"`           // 0默认,1管理员
}

var crud = gt.NewCrud(
	gt.Model(Admin{}),
)

// get data, by id
func (c *Admin) GetByID(id string) interface{} {

	var data Admin // not use *Admin
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Admin) GetBySearch(params map[string][]string) interface{} {
	var datas []*Admin
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *Admin) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *Admin) Update(data *Admin) interface{} {

	data.Password = util.AesEn(data.Password)
	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *Admin) Create(data *Admin) interface{} {

	// create time
	//(*data).Createtime = time2.CTime(time.Now())
	data.Password = util.AesEn(data.Password)

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
