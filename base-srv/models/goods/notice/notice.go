package notice

import (
	"demo/base-srv/controllers/wx/template_msg"
	"demo/commons/models"
	"github.com/dreamlu/gt"
)

// 社群特有字段
// 商品--开团预告模型
type GsNotice struct {
	models.AdminCom
	ClientID uint64 `json:"client_id" gorm:"type:bigint(20);UNIQUE_INDEX:开团提醒已添加"`
	GoodsID  uint64 `json:"goods_id" gorm:"type:bigint(20);UNIQUE_INDEX:开团提醒已添加"`
	//StartTime time.CTime `json:"start_time" gorm:"type:datetime"` // 开团时间
	template_msg.ModelMsg
	Status byte `json:"status" gorm:"type:tinyint(1);DEFAULT:0"` // 0默认未通知, 1已通知
}

var crud = gt.NewCrud(
	gt.Model(GsNotice{}),
)

// get data, limit and search
// ClientPage 1, everyPage 10 default
//func (c *GsNotice) Search(params cmap.CMap) (datas []*GsNotice, pager result.Pager, err error) {
//
//	crud.Params(
//		gt.Data(&datas),
//	)
//	cd := crud.GetBySearch(params)
//	if cd.Error() != nil {
//		return nil, pager, cd.Error()
//	}
//	return datas, cd.Pager(), nil
//}

// create data
func (c *GsNotice) Create(data *GsNotice) (*GsNotice, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}
