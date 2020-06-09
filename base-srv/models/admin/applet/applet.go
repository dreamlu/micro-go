package applet

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	models2 "micro-go/commons/models"
	"strings"
)

// 小程序相关 模型
type Applet struct {
	models2.AdminCom
	Appid     string `json:"appid" gorm:"type:varchar(50);UNIQUE_INDEX:appid已存在"` // appid
	Secret    string `json:"secret" gorm:"type:varchar(50)"`                      // secret
	MchID     string `json:"mch_id" gorm:"type:varchar(50)"`                      // 商户号
	PaySecret string `json:"pay_secret" gorm:"type:varchar(50)"`                  // 商户秘钥
	// 证书
	AppCert string `json:"app_cert"` // apiclient_cert.pem
	AppKey  string `json:"app_key"`  // apiclient_key.pem
	// logo
	Logo string `json:"logo"` // logo
}

var crud = gt.NewCrud(
	gt.Model(Applet{}),
)

// get data, by id
func (c *Applet) GetByAdminID(admin_id uint64) error {

	//var data Applet // not use *Applet
	sql := fmt.Sprintf("select %s from applet where admin_id = ?", gt.GetColSQL(Applet{}))
	cd := crud.Params(gt.Data(&c)).Select(sql, admin_id).Single()
	if err := cd.Error(); err != nil {
		//log.Log.Error(err.Error())
		return err
	}
	//c = &data
	return nil
}

// get data, by id
func (c *Applet) GetByAppid(appid string) error {

	//var data Applet // not use *Applet
	sql := fmt.Sprintf("select %s from applet where appid = ?", gt.GetColSQL(Applet{}))
	cd := crud.Params(gt.Data(c)).Select(sql, appid).Single()
	if err := cd.Error(); err != nil {
		//log.Log.Error(err.Error())
		return err
	}
	//c = &data
	return nil
}

// get data, by id
func (c *Applet) GetByID(id string) interface{} {

	var data Applet // not use *Applet
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// GsCategoryPage 1, everyPage 10 default
func (c *Applet) GetBySearch(params cmap.CMap) (datas []*Applet, pager result.Pager, err error) {

	crud.Params(
		gt.Data(&datas),
	)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// delete data, by id
func (c *Applet) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *Applet) Update(data *Applet) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *Applet) Create(data *Applet) interface{} {

	data.Appid = strings.Trim(data.Appid, "")
	data.Secret = strings.Trim(data.Secret, "")
	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
