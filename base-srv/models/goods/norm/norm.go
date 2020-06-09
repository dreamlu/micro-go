package norm

import (
	"github.com/dreamlu/gt"
	models2 "micro-go/commons/models"
)

// 商品规格值
type GsNorm struct {
	models2.ModelCom
	GsNormPar
}

type GsNormPar struct {
	GoodsID    uint64  `json:"goods_id" gorm:"type:bigint(20)"`          // 商品id
	NameFirst  string  `json:"name_first" gorm:"type:varchar(20)"`       // 一级规格
	NameSecond string  `json:"name_second" gorm:"type:varchar(20)"`      // 二级规格
	NameThird  string  `json:"name_third" gorm:"type:varchar(20)"`       // 三级规格, 最多三层规格, 这样设计比用json简单, 容易对比取值计算
	Price      float64 `json:"price" gorm:"type:double(10,2);DEFAULT:0"` // 价格
	Inventory  int     `json:"inventory" gorm:"type:int(11)"`            // 库存
	Img        string  `json:"img"`                                      // 规格图
	Code       string  `json:"code" gorm:"type:varchar(50)"`             // 商品编码
	SellNum    int     `json:"sell_num" gorm:"type:int(11);DEFAULT:0"`   // 已卖数量

	// 社群特有
	HeadPrice   float64 `json:"head_price" gorm:"type:double(10,2)"`    // 团长价
	SpHeadPrice float64 `json:"sp_head_price" gorm:"type:double(10,2)"` // 超级团长价
	//SpuCode     string  `json:"spu_code" gorm:"type:varchar(50)"`       // spu编码
}

// 规格嵌套
type GsNormSt struct {
	Name  string      `json:"name"`
	Child interface{} `json:"child"`
}

// 最后一层数据
type GsNormStCh struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price" gorm:"type:double(10,2);DEFAULT:0"` // 价格
	Inventory int     `json:"inventory" gorm:"type:int(11)"`            // 库存
	Img       string  `json:"img"`                                      // 规格图
}

var crud = gt.NewCrud(
	gt.Model(GsNorm{}),
)

// 批量创建和修改
func CreateAndUpdate(cd gt.Crud, goods_id uint64, data []*GsNormPar) error {

	if cd != nil {
		crud = cd
	}
	// 删除
	if goods_id != 0 {
		err := crud.Select("delete from gs_norm where goods_id = ?", goods_id).Exec().Error()
		if err != nil {
			gt.Logger().Error(err.Error())
			// 继续执行
		}
	}
	if len(data) == 0 {
		return nil
	}

	// 重新创建
	crud.Params(
		gt.Table("gs_norm"),
		gt.Model(GsNormPar{}),
		gt.Data(data),
	)
	if err := crud.CreateMore().Error(); err != nil {
		return err
	}
	return nil
}
