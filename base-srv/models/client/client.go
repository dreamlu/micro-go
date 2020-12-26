package client

import (
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
)

// 客户模型
type Client struct {
	models2.ModelCom
	AdminID uint64  `json:"admin_id" gorm:"type:bigint(20);INDEX:查询索引admin_id;UNIQUE_INDEX:openid已存在"`
	Name    string  `gorm:"type:varchar(30)" json:"name" valid:"required,len=2-20"` // 昵称
	Openid  string  `json:"openid" gorm:"varchar(30);UNIQUE_INDEX:openid已存在"`       // openID
	HeadImg string  `json:"head_img"`                                               // 头像
	Phone   *string `json:"phone" gorm:"type:varchar(20);UNIQUE_INDEX:手机号已存在"`      // 手机号
}

type ClientDP struct {
	Client        // 头像
	Type   *byte  `json:"type"`     // 0团长, 1超级团长, 2团员
	DisID  uint64 `json:"dis_id"`   // 分销人id
	PName  string `json:"p_name"`   // 上级分销人名称
	PDisID uint64 `json:"p_dis_id"` // 上级分销人id
}

type ClientD struct {
	ClientDP
	BuyNum int8    `json:"buy_num" gt:"sub_sql"` // 购买次数
	Money  float64 `json:"money" gt:"sub_sql"`   // 实付金额
}

// 个人详情中心
type ClientC struct {
	ClientDP
	Status []*Status `json:"status"`
}

type Status struct {
	Status int64 `json:"status"`
	Num    int64 `json:"num"`
}

var crud = gt.NewCrud(
	gt.Model(Client{}),
)

// get data, by id
func (c *Client) GetByIDC(params cmap.CMap) (data ClientC, err error) {

	cld, er := c.GetByID(params)
	if er != nil {
		return data, nil
	}
	// 查找客户订单中的各个状态
	data.ClientDP = cld
	crud.Params(gt.Data(&data.Status)).
		Select("select status,count(*) num from `order` where client_id = ? and status in (0,1,2,3,5) group by status", params.Get("id")).
		Single()

	return
}

// get data, by id
func (c *Client) GetByID(params cmap.CMap) (data ClientDP, err error) {

	typeSQL := "(select type from dis.dis where client_id = `client`.id) as type"
	disIDSQL := "(select id from dis.dis where client_id = `client`.id and admin_id = `client`.admin_id) as dis_id"
	pnameSQL := "(select cc.name from dis.dis_client aa inner join dis.dis bb on aa.dis_id = bb.id inner join client cc on bb.client_id = cc.id where aa.client_id = `client`.id) as p_name"
	pidSQL := "(select bb.id from dis.dis_client aa inner join dis.dis bb on aa.dis_id = bb.id where aa.client_id = `client`.id) as p_dis_id"
	crud.Params(gt.Data(&data), gt.SubSQL(typeSQL, disIDSQL, pnameSQL, pidSQL))
	if err = crud.GetByData(params).Error(); err != nil {
		return
	}
	if data.Type == nil {
		ty := byte(2)
		data.Type = &ty // 团员
	}
	return
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *Client) GetBySearch(params cmap.CMap) interface{} {

	buyNumSQL := "(select count(*) from `order` where client_id = `client`.id and status in (1,2,3,7)) as buy_num"
	moneySQL := "(select sum(money) from `order` where client_id = `client`.id and status in (1,2,3,7)) as money"
	var datas []*ClientD
	crud.Params(
		gt.Data(&datas),
		gt.SubSQL(buyNumSQL, moneySQL),
	)
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return result.CError(cd.Error())
	}
	return result.GetSuccessPager(datas, cd.Pager())
}

// delete data, by id
func (c *Client) Delete(id string) interface{} {

	if err := crud.Delete(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update data
func (c *Client) Update(data *Client) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *Client) Create(data *Client) interface{} {

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		//log.Log.Error(err.Error())
		return result.CError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate).Add("id", data.ID)
}
