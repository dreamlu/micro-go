package comment

import (
	models2 "demo/commons/models"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/json"
)

// 订单--评论模型
type OrderComment struct {
	models2.ModelCom
	OrderCommentPar
}

type OrderCommentPar struct {
	GoodsID  uint64     `json:"goods_id" gorm:"type:bigint(20)"`       // 商品id
	OrderID  uint64     `json:"order_id" gorm:"type:bigint(20)"`       // 订单id
	ClientID uint64     `json:"client_id" gorm:"type:bigint(20)"`      // 客户id
	Star     *int8      `json:"star" gorm:"type:tinyint(2);DEFAULT:5"` // 星级, 0-5颗星,默认5颗星
	Content  string     `json:"content"`                               // 评价内容
	Img      json.CJSON `json:"img" gorm:"type:json"`                  // 评价图
}

type OrderCommentD struct {
	OrderComment
	ClientName    string `json:"client_name"`
	ClientHeadImg string `json:"client_head_img"`
}

var crud = gt.NewCrud(
	gt.Model(OrderComment{}),
)

// get data, by id
func (c *OrderComment) GetByID(id string) (*OrderComment, error) {

	var data OrderComment // not use *OrderComment
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
// OrderCommentPage 1, everyPage 10 default
func (c *OrderComment) GetBySearch(params cmap.CMap) (datas []*OrderCommentD, pager result.Pager, err error) {

	crud.Params(
		gt.Data(&datas),
		gt.Model(OrderCommentD{}),
		gt.InnerTable([]string{"order_comment", "client"}),
	)
	cd := crud.GetMoreBySearch(params)
	if cd.Error() != nil {
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// delete data, by id
func (c *OrderComment) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *OrderComment) Update(data *OrderComment) (*OrderComment, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *OrderComment) Create(data *OrderComment) (*OrderComment, error) {

	cd := gt.NewCrud().Begin()
	cd.Params(gt.Data(data))
	if err := cd.Create().Error(); err != nil {
		cd.Rollback()
		return nil, err
	}

	// 更改订单状态--已评价
	cd.Select("update `order` set `status` = 7 where id = ?", data.OrderID).Exec()
	if err := cd.Error(); err != nil {
		cd.Rollback()
		return nil, err
	}

	if err := cd.Commit().Error(); err != nil {
		cd.Rollback()
		return nil, err
	}

	return data, nil
}

// create data
func (c *OrderComment) CreateMore(data []*OrderCommentPar) error {

	cd := gt.NewCrud().Begin()
	cd.Params(gt.Table("order_comment"), gt.Model(OrderCommentPar{}), gt.Data(data))
	if err := cd.CreateMore().Error(); err != nil {
		cd.Rollback()
		return err
	}

	// 更改订单状态--已评价
	cd.Select("update `order` set `status` = 7 where id = ?", data[0].OrderID).Exec()
	if err := cd.Error(); err != nil {
		cd.Rollback()
		return err
	}

	if err := cd.Commit().Error(); err != nil {
		cd.Rollback()
		return err
	}

	return nil
}
