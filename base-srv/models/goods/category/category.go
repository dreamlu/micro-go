package category

import (
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
)

// 商品--分类模型
type GsCategory struct {
	models2.ModelCom
	AdminID uint64 `json:"admin_id" gorm:"type:bigint(20);UNIQUE_INDEX:分类已存在"`
	Name    string `gorm:"type:varchar(30);UNIQUE_INDEX:分类已存在" json:"name"` // 名称
	IsShow  *int8  `json:"is_show" gorm:"tinyint(2);DEFAULT:1"`             // 0不显示,1显示(默认)
	Pid     uint64 `json:"pid" gorm:"type:bigint(20);DEFAULT:0"`            // 父级id
	Icon    string `json:"icon"`                                            // 图标
}

// 二级列表
type GsCategoryList struct {
	ID    uint64            `json:"id"`
	Name  string            `json:"name"` // 名称
	Icon  string            `json:"icon"` // 图标
	Child []*GsCategoryList `json:"child"`
}

var crud = gt.NewCrud(
	gt.Model(GsCategory{}),
)

// get data, limit and search
// GsCategoryPage 1, everyPage 10 default
func (c *GsCategory) List(params cmap.CMap) (datas []*GsCategoryList, err error) {

	// 全部数据
	params.Add("is_show", "1")
	cas, _, err := c.GetBySearch(params)
	if err != nil {
		return nil, err
	}
	var data *GsCategoryList
	// 格式转换
	for _, v := range cas {
		if v.Pid == 0 {
			data = &GsCategoryList{
				ID:   v.ID,
				Name: v.Name,
				Icon: v.Icon,
			}
			datas = append(datas, data)
			//cas = append(cas[:k], cas[:k+1]...)
		}
	}
	for _, d := range datas {
		for _, v := range cas {
			if d.ID == v.Pid {
				data = &GsCategoryList{
					ID:   v.ID,
					Name: v.Name,
					Icon: v.Icon,
				}
				d.Child = append(d.Child, data)
			}
		}
	}
	// 转换完成
	return datas, nil
}

// get data, by id
func (c *GsCategory) GetByID(id string) (*GsCategory, error) {

	var data GsCategory // not use *GsCategory
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
// GsCategoryPage 1, everyPage 10 default
func (c *GsCategory) GetBySearch(params cmap.CMap) (datas []*GsCategory, pager result.Pager, err error) {

	// params
	// p: 0筛选以及分类, 1筛选二级分类
	p := params.Get("p")
	var parWhereSQL string
	if p == "0" {
		parWhereSQL = "pid = 0"
		params.Del("p")
	}
	if p == "1" {
		parWhereSQL = "pid != 0"
		params.Del("p")
	}
	cd := gt.NewCrud(
		gt.Model(GsCategory{}),
		gt.Data(&datas),
		gt.SubWhereSQL(parWhereSQL),
	)
	cd = cd.GetBySearch(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// delete data, by id
func (c *GsCategory) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *GsCategory) Update(data *GsCategory) (*GsCategory, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *GsCategory) Create(data *GsCategory) (*GsCategory, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}
