package cart

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/json"
	"micro-go/base-srv/models/goods"
	models2 "micro-go/commons/models"
	"strings"
)

// 购物车模型
type Cart struct {
	models2.AdminCom
	ClientID uint       `json:"client_id" gorm:"type:bigint(20)"` // 客户id
	GoodsID  uint       `json:"goods_id" gorm:"type:bigint(20)"`  // 商品id
	Name     string     `gorm:"type:varchar(50)" json:"name"`     // 商品名称
	ImgUrl   string     `json:"img_url"`                          // 封面
	Price    float64    `json:"price" gorm:"type:double(10,2)"`   // 价格
	Num      int        `json:"num" gorm:"type:int(11)"`          // 数量
	CNorm    json.CJSON `json:"c_norm" gorm:"type:json"`          // 客户规格,json
}

// 购物车详情
type CartDe struct {
	Cart
	Goods  goods.GoodsMD `json:"goods"`
	IsExit int8          `json:"is_exit" gt:"sub_sql"` // 是否存在,0否,1是
}

var crud = gt.NewCrud(
	gt.Model(Cart{}),
)

// get data, by id
func (c *Cart) GetByID(id string) interface{} {

	var data Cart // not use *Cart
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetSuccess(data)
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Cart) GetBySearch(params map[string][]string) interface{} {

	isExitSQL := "(select 1 from goods a where a.id = cart.goods_id and a.is_shelf = 1) as is_exit"
	var datas []*CartDe
	crud.Params(
		gt.Data(&datas),
		gt.SubSQL(isExitSQL),
	)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		return result.CError(cd.Error())
	}

	// 查找商品详情
	for _, v := range datas {
		if v.IsExit == 1 {
			gs, err := v.Goods.Goods.Get(v.GoodsID)
			if err != nil {
				gt.Logger().Error("[购物车查找商品详情错误]", err.Error())
				// 继续执行下一个
			}
			v.Goods = gs
		}
	}

	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *Cart) Delete(id string) interface{} {

	strings.Replace(id, "'", "\\'", -1)
	if err := crud.Select("delete from `cart` where id in (" + id + ")").Exec().Error(); err != nil {
		return result.GetError(err)
	}
	return result.MapDelete
}

// update data
func (c *Cart) Update(data *Cart) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *Cart) Create(data *Cart) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
