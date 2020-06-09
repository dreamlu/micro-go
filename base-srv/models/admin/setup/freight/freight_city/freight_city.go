package freight_city

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	models2 "micro-go/commons/models"
	"strings"
)

// 运费模板 模型
type FreightCity struct {
	models2.ModelCom
	FreightCityPar
}

type FreightCityPar struct {
	FreightID uint64  `json:"freight_id" gorm:"type:bigint(20);UNIQUE_INDEX:城市已存在"`  // 模板id
	Code      int     `json:"code" gorm:"type:int(11);DEFAULT:0;UNIQUE_INDEX:城市已存在"` // 城市编码code, 全国默认0
	City      string  `json:"city" gorm:"type:varchar(30)"`                          // 城市
	Price     float64 `json:"price" gorm:"type:double(10,2);DEFAULT:0"`              // 运费价格
}

var crud = gt.NewCrud(
	gt.Model(FreightCity{}),
)

// get data, by id
func (c *FreightCity) Get(id string) (*FreightCity, error) {

	var data FreightCity // not use *Freight
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

func search(params cmap.CMap) (datas []*FreightCity, pager result.Pager, err error) {

	crud.Params(
		gt.Model(FreightCity{}),
		gt.Data(&datas),
	)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}

	// 没数据, 默认返回全国
	if len(datas) == 0 && params.Get("freight_id") != "" && params.Get("city") != "" {
		params.Del("city")
		params.Add("code", "0")
		return search(params)
	}

	return datas, cd.Pager(), nil
}

// get data, limit and search
// FreightPage 1, everyPage 10 default
func (c *FreightCity) Search(params cmap.CMap) (datas []*FreightCity, pager result.Pager, err error) {

	fids := strings.Split(params.Get("freight_id"), ",")
	//city := params.Get("city")

	if len(fids) == 1 {
		return search(params)
	}
	// no pager
	for _, v := range fids {
		params.Del("freight_id")
		params.Add("freight_id", v)
		//params.Add("city",city)
		datas2, _, err := search(params)
		if err != nil {
			return datas, pager, err
		}
		datas = append(datas, datas2[:]...)
	}
	return
}

// delete data, by id
func (c *FreightCity) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *FreightCity) Update(data *FreightCity) (*FreightCity, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *FreightCity) Create(data *FreightCity) (*FreightCity, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *FreightCity) CreateAll(data []*FreightCityPar) ([]*FreightCityPar, error) {

	crud.Params(
		gt.Model(FreightCityPar{}),
		gt.Table("freight_city"),
		gt.Data(data),
	)
	if err := crud.CreateMoreData().Error(); err != nil {
		return nil, err
	}
	return data, nil
}
