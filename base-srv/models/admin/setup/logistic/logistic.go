package logistic

import (
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
)

// 物流公司模型
type Logistic struct {
	models2.AdminCom
	Name string `json:"name" gorm:"type:varchar(30)"`
	Com  string `json:"com" gorm:"type:varchar(30)"` // 快递公司编码
}

var crud = gt.NewCrud(
	gt.Model(Logistic{}),
)

// get data, by id
func (c *Logistic) GetByName(name string) error {

	var (
		params = cmap.CMap{}
	)
	params.Add("name", name)
	cd := gt.NewCrud(gt.Model(Logistic{}), gt.Data(&c))
	if err := cd.GetByData(params).Error(); err != nil {
		//log.Log.Error(err.Error())
		return err
	}
	return nil
}

// get data, by id
func (c *Logistic) GetByID(id string) (*Logistic, error) {

	var data Logistic // not use *Logistic
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
// LogisticPage 1, everyPage 10 default
func (c *Logistic) GetBySearch(params cmap.CMap) (datas []*Logistic, pager result.Pager, err error) {

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
func (c *Logistic) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *Logistic) Update(data *Logistic) (*Logistic, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *Logistic) Create(data *Logistic) (*Logistic, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}
