package record

import (
	models2 "demo/commons/models"
	json2 "encoding/json"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/log"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
)

// TODO 删除过久日志
// 记录,操作日志
type Record struct {
	models2.AdminCom
	Type int8       `json:"type" gorm:"type:tinyint(2)"` // 0(订单)
	Info string     `json:"info"`                        // 操作详情
	Data json.CJSON `json:"data" gorm:"type:json"`       // 操作的字段数据,可不做展示
}

type RecordD struct {
	models2.ModelCom
	Type int8   `json:"type"`
	Info string `json:"info"`
}

var crud = gt.NewCrud(
	gt.Model(Record{}),
)

// get data, limit and search
// ClientPage 1, everyPage 10 default
// TODO 增加时间段筛选
func (c *Record) Search(params cmap.CMap) (datas []*RecordD, pager result.Pager, err error) {

	crud.Params(
		gt.Data(&datas),
	)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// create data
func (c *Record) Create(data *Record) (*Record, error) {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

// 日志记录
func Log(ty int8, admin_id uint64, args ...interface{}) {
	var re = Record{Type: ty, AdminCom: models2.AdminCom{AdminID: admin_id}}
	if len(args) > 0 {
		re.Info = args[0].(string)
	}
	if len(args) > 1 {
		re.Data, _ = json2.Marshal(args[1])
	}
	_, err := re.Create(&re)
	if err != nil {
		log.Error("[日志记录错误]:", err.Error())
	}
}
