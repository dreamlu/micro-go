package goods

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
	"github.com/dreamlu/gt/tool/type/time"
	"micro-go/base-srv/models/goods/norm"
	"micro-go/base-srv/util/cron"
	"micro-go/commons/models"
)

// 商品模型
type Goods struct {
	models.AdminCom
	Name         string     `gorm:"type:varchar(30)" json:"name" valid:"required,len=2-20"` // 名称
	Price        float64    `json:"price" gorm:"type:double(10,2)"`                         // 现价
	OldPrice     float64    `json:"old_price" gorm:"type:double(10,2)"`                     // 原价
	ImgUrl       string     `json:"img_url"`                                                // 封面
	Carousel     json.CJSON `json:"carousel" gorm:"type:json"`                              // 轮播图
	Introduce    string     `json:"introduce" gorm:"type:longtext"`                         // 介绍,富文本
	IsShelf      *int8      `json:"is_shelf" gorm:"type:tinyint(2);DEFAULT:1"`              // 是否上架,0否,1上架(默认)
	Inventory    int        `json:"inventory" gorm:"type:int(11)" gt:"sub_sql"`             // 库存
	SellNum      int        `json:"sell_num" gorm:"type:int(11);DEFAULT:0" gt:"sub_sql"`    // 已卖数量
	GsCategoryID string     `json:"gs_category_id" gorm:"type:varchar(50)"`                 // 分类id(一二级都支持), 逗号分割
	NormInfo     json.CJSON `json:"norm_info" gorm:"type:json"`                             // 规格对应的设置, 名称设置, 值的设置参考
	FreightID    uint64     `json:"freight_id" gorm:"type:bigint(20)"`                      // 模板id
	IsRec        *byte      `json:"is_rec" gorm:"type:tinyint(2);DEFAULT:0"`                // 是否首页推荐, 0否,1是
	Code         string     `json:"code" gorm:"type:varchar(50)"`                           // 商品编码
	VirtualNum   uint64     `json:"virtual_num" gorm:"type:int(11)"`                        // 虚拟销量

	// update
	HeadPrice   float64 `json:"head_price" gorm:"type:double(10,2)"`    // 团长价
	SpHeadPrice float64 `json:"sp_head_price" gorm:"type:double(10,2)"` // 超级团长价

	// 社群特有字段
	Share       json.CJSON `json:"share" gorm:"type:json"`              // 分享素材,json
	StartTime   time.CTime `json:"start_time" gorm:"type:datetime"`     // 开始时间
	EndTime     time.CTime `json:"end_time" gorm:"type:datetime"`       // 结束时间
	Label       json.CJSON `json:"label" gorm:"type:json"`              // 标签,json
	WarehouseID uint64     `json:"warehouse_id" gorm:"type:bigint(20)"` // 货仓id
	SpuCode     string     `json:"spu_code" gorm:"type:varchar(50)"`    // spu编码, 商品唯一, 和多规格无关
	IsG         int8       `json:"is_g"`                                // 是否官方,0否,1是
}

type GoodsD struct {
	Goods
	WarehouseName string `json:"warehouse_name" gt:"field:gs_warehouse.name"`
}

// 商品详情
type GoodsM struct {
	Goods
	Norm []*norm.GsNormPar `json:"norm"`
}

// 商品详情
type GoodsMD struct {
	Goods
	Norm []*norm.GsNorm `json:"norm"`
}

var crud = gt.NewCrud(
	gt.Model(Goods{}),
)

// 已卖数量加上待付款数量
var sellNumSQL = "((select ifnull(sum(bb.num),0) from `order` aa inner join `order_goods` bb on aa.id = bb.order_id where aa.status = 0 and bb.goods_id = `goods`.id) + `goods`.sell_num + (select ifnull(sum(bb.sell_num),0) from `gs_norm` bb where bb.goods_id = `goods`.id)) sell_num"

// 总库存
var inventorySQL = "(ifnull((select sum(inventory) from gs_norm where goods_id = `goods`.id),0) + ifnull(`goods`.inventory,0)) as inventory"

// get data, by id
func (c *Goods) GetTrans(cd gt.Crud, id interface{}) (data Goods, err error) {

	cd.Params(gt.Data(&data))
	if err = cd.GetByID(id).Error(); err != nil {
		return
	}
	return
}

// get data, by id
func (c *Goods) GetByID(params cmap.CMap) (data GoodsMD, err error) {

	crud.Params(
		gt.Data(&data.Goods),
		gt.SubSQL(sellNumSQL, inventorySQL),
	)
	if err = crud.GetByData(params).Error(); err != nil {
		return
	}

	// 查找规格
	crud.Params(gt.Data(&data.Norm)).Select(fmt.Sprintf("select %s from `gs_norm` where goods_id = ?", gt.GetColSQL(norm.GsNorm{})), data.ID).Single()
	if err = crud.Error(); err != nil {
		return
	}
	return
}

// get data, by id
func (c *Goods) Get(id interface{}) (data GoodsMD, err error) {

	crud.Params(
		gt.Data(&data.Goods),
		gt.SubSQL(sellNumSQL, inventorySQL),
	)
	if err = crud.GetByID(id).Error(); err != nil {
		return
	}

	// 查找规格
	crud.Params(gt.Data(&data.Norm)).Select(fmt.Sprintf("select %s from `gs_norm` where goods_id = ?", gt.GetColSQL(norm.GsNorm{})), data.ID).Single()
	if err = crud.Error(); err != nil {
		return
	}
	return
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Goods) Search(params cmap.CMap) interface{} {

	// 上架
	//params.Add("is_shelf", "1")
	// 分类匹配
	// 字段中逗号分割存储(ps:需要计算的值不要用这种形式存储,单独建表)
	var gcidWhereSQL string
	if gcid := params.Get("gs_category_id"); gcid != "" {
		gcidWhereSQL = "gs_category_id like '%" + gcid + "%'"
		params.Del("gs_category_id")
	}
	// 社群特有字段
	// 时间段,status, 0预告,1进行中,2结束
	statusSQL := ""
	if status := params.Get("status"); status != "" {
		switch status {
		case "0":
			statusSQL = "now() < `goods`.start_time"
		case "1":
			statusSQL = "now() <= `goods`.end_time and now() >= `goods`.start_time"
		case "2":
			statusSQL = "now() > `goods`.end_time"
		}
		params.Del("status")
	}

	var datas []*GoodsD
	crud2 := gt.NewCrud(
		gt.Model(GoodsD{}),
		gt.Data(&datas),
		gt.SubWhereSQL(gcidWhereSQL, statusSQL),
		gt.SubSQL(inventorySQL, sellNumSQL),
		gt.LeftTable([]string{"goods:warehouse_id", "gs_warehouse:id"}),
	)
	cd := crud2.GetMoreBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *Goods) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *Goods) Update(data *GoodsM) interface{} {

	cd := gt.NewCrud().Begin()

	exitData, err := c.GetTrans(cd, data.ID)
	if err != nil {
		cd.Rollback()
		return result.CError(err)
	}

	cd.Params(gt.Data(data.Goods))
	if err := cd.Update().Error(); err != nil {
		cd.Rollback()
		return result.CError(err)
	}

	// 修改相关规格
	if data.Norm != nil {
		err := norm.CreateAndUpdate(cd, data.ID, data.Norm)
		if err != nil {
			cd.Rollback()
			return result.CError(err)
		}
	}

	if err := cd.Commit().Error(); err != nil {
		cd.Rollback()
		return result.CError(err)
	}

	go func() {
		t := time.CTime{}
		if data.StartTime != t &&
			data.StartTime != exitData.StartTime {
			// 进行定时检测通知任务的创建
			cron.CronNotice(data.StartTime)
		}
	}()

	return result.MapUpdate
}

// create data
func (c *Goods) Create(data *GoodsM) interface{} {

	//gt.Logger().Info(data.StartTime.String())
	cd := gt.NewCrud().Begin()
	cd.Params(gt.Data(&data.Goods))
	if err := cd.Create().Error(); err != nil {
		cd.Rollback()
		return result.CError(err)
	}

	// 创建相关规格
	for _, v := range data.Norm {
		v.GoodsID = data.ID
	}
	if len(data.Norm) > 0 {
		err := norm.CreateAndUpdate(cd, 0, data.Norm)
		if err != nil {
			cd.Rollback()
			return result.CError(err)
		}
	}

	if err := cd.Commit().Error(); err != nil {
		cd.Rollback()
		return result.CError(err)
	}

	// 进行定时检测通知任务的创建
	go cron.CronNotice(data.StartTime)

	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
