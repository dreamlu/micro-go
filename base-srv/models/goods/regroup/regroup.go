package regroup

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/time"
	"micro-go/commons/models"
)

// 社群特有字段
// 商品--复团申请
type GsRegroup struct {
	models.ModelCom
	GoodsID   uint64     `json:"goods_id" gorm:"type:bigint(20);UNIQUE_INDEX:复团已申请"`
	ClientID  uint64     `json:"client_id" gorm:"type:bigint(20);UNIQUE_INDEX:复团已申请"`
	StartTime time.CTime `json:"start_time" gorm:"type:datetime;UNIQUE_INDEX:复团已申请"` // 开始时间
	EndTime   time.CTime `json:"end_time" gorm:"type:datetime;UNIQUE_INDEX:复团已申请"`   // 结束时间
}

type GsRegroupCount struct {
	StartTime time.CTime `json:"start_time"` // 开始时间
	EndTime   time.CTime `json:"end_time"`   // 结束时间
	Num       int64      `json:"num"`        // 复团数量
}

var crud = gt.NewCrud(
	gt.Model(GsRegroup{}),
)

func (c *GsRegroup) Count(params cmap.CMap) (datas []*GsRegroupCount, pager result.Pager, err error) {

	crud.Params(
		gt.Data(&datas),
	)
	cd := crud.Select("select start_time,end_time,count(*) as num "+
		"from gs_regroup "+
		"where goods_id = ? "+
		"group by goods_id,start_time,end_time "+
		"order by start_time desc", params.Get("goods_id"))
	cd.Search(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// get data, by id
func (c *GsRegroup) GetByID(id string) (*GsRegroup, error) {

	var data GsRegroup // not use *GsRegroup
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
func (c *GsRegroup) GetBySearch(params cmap.CMap) (datas []*GsRegroup, pager result.Pager, err error) {

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
func (c *GsRegroup) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *GsRegroup) Update(data *GsRegroup) (*GsRegroup, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *GsRegroup) Create(data *GsRegroup) (*GsRegroup, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}
