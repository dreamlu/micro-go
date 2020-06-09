package warehouse

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"micro-go/commons/models"
)

// 社群特有字段
// 商品--货仓
type GsWarehouse struct {
	models.AdminCom
	Name string `json:"name" gorm:"type:varchar(50)"`
}

var crud = gt.NewCrud(
	gt.Model(GsWarehouse{}),
)

// get data, by id
func (c *GsWarehouse) GetByID(id string) (*GsWarehouse, error) {

	var data GsWarehouse // not use *GsWarehouse
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
func (c *GsWarehouse) GetBySearch(params cmap.CMap) (datas []*GsWarehouse, pager result.Pager, err error) {

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
func (c *GsWarehouse) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *GsWarehouse) Update(data *GsWarehouse) (*GsWarehouse, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *GsWarehouse) Create(data *GsWarehouse) (*GsWarehouse, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}
