package freight

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"micro-go/base-srv/models/admin/setup/freight/freight_city"
	models2 "micro-go/commons/models"
)

// 运费模板 模型
type Freight struct {
	models2.AdminCom
	Name string `json:"name" gorm:"type:varchar(20)"` // 名称
}

var crud = gt.NewCrud(
	gt.Model(Freight{}),
)

// get data, by id
func (c *Freight) Get(id string) (*Freight, error) {

	var data Freight // not use *Freight
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
// FreightPage 1, everyPage 10 default
func (c *Freight) Search(params cmap.CMap) (datas []*Freight, pager result.Pager, err error) {

	adminIDSQL := ""
	if admin_id := params.Get("admin_id"); admin_id != "" {
		adminIDSQL = "admin_id in (" + admin_id + ")"
		params.Del("admin_id")
	}
	crud.Params(
		gt.Data(&datas),
		gt.SubWhereSQL(adminIDSQL),
	)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// delete data, by id
func (c *Freight) Delete(id string) error {

	err := crud.Delete(id).Error()
	if err == nil {
		// 删除关联表信息
		go func() {
			crud.Select("delete from freight_city where freight_id = ?", id).Exec()
		}()
	}

	return err
}

// update data
func (c *Freight) Update(data *Freight) (*Freight, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *Freight) Create(data *Freight) (*Freight, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}

	// 默认创建一个全国的
	var fc = freight_city.FreightCity{
		FreightCityPar: freight_city.FreightCityPar{
			FreightID: data.ID,
			City:      "默认全国",
		},
	}
	_, err := fc.Create(&fc)
	if err != nil {
		return nil, err
	}

	return data, nil
}
