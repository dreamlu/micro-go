package address

import (
	"demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
)

// 客户地址
type ClientAddress struct {
	models.ModelCom
	ClientAddressP
}

type ClientAddressP struct {
	ClientID uint64 `json:"client_id"`                        // 客户id
	Name     string `json:"name" gorm:"type:varchar(20)"`     // 姓名
	Sex      int8   `json:"sex" gorm:"default:0"`             // 0未知, 1男, 2女
	Phone    string `json:"phone" gorm:"type:varchar(20)"`    // 手机号
	Address  string `json:"address"`                          // 详细地址
	HouseNum string `json:"house_num"`                        // 门牌号
	Province string `json:"province" gorm:"type:varchar(50)"` // 省份
	City     string `json:"city" gorm:"type:varchar(50)"`     // 城市
	Area     string `json:"area" gorm:"type:varchar(50)"`     // 区
}

var crud = gt.NewCrud(
	gt.Model(ClientAddress{}),
)

// get data, by id
func (c *ClientAddress) Get(params cmap.CMap) (data ClientAddress, err error) {
	crud.Params(gt.Data(&data))
	if err = crud.GetByData(params).Error(); err != nil {
		return
	}
	return
}

// get data, limit and search
func (c *ClientAddress) Search(params cmap.CMap) (datas []*ClientAddress, pager result.Pager, err error) {

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
func (c *ClientAddress) Delete(id interface{}) error {

	return crud.Delete(id).Error()
}

// update data
func (c *ClientAddress) Update(data *ClientAddress) (*ClientAddress, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

// create data
func (c *ClientAddress) Create(data *ClientAddress) (*ClientAddress, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

// create data
func (c *ClientAddress) CreateMore(data []*ClientAddressP) error {

	cd := gt.NewCrud(gt.Data(data), gt.Table("client_address"), gt.Model(ClientAddressP{}))
	cd.Params(gt.Data(data))
	if err := cd.CreateMore().Error(); err != nil {
		return err
	}
	return nil
}
