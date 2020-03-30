// @author  dreamlu
package models

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"micro-go/commons/models"
)

/*user model*/
type User struct {
	models.ModelCom
	Name string `json:"name" valid:"required,len=2-20"`
}

var crud = gt.NewCrud(
	gt.Model(User{}),
)

// get data, by id
func (c *User) GetByID(id string) (*User, error) {

	var data User // not use *User
	crud.Params(gt.Data(&data))
	if err := crud.GetByID(id).Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return &data, nil
}

// get data, limit and search
// clientPage 1, everyPage 10 default
func (c *User) GetBySearch(params cmap.CMap) (datas []*User, pager result.Pager, err error) {
	//var datas []*User
	crud.Params(gt.Data(&datas))
	cd := crud.GetBySearch(params)
	if cd.Error() != nil {
		//log.Log.Error(err.Error())
		return nil, pager, cd.Error()
	}
	return datas, cd.Pager(), nil
}

// delete data, by id
func (c *User) Delete(id string) error {

	return crud.Delete(id).Error()
}

// update data
func (c *User) Update(data *User) (*User, error) {

	crud.Params(gt.Data(data))
	if err := crud.Update().Error(); err != nil {
		//log.Log.Error(err.Error())
		return nil, err
	}
	return data, nil
}

// create data
func (c *User) Create(data *User) (*User, error) {

	// create time
	//(*data).Createtime = time2.CTime(time.Now())

	crud.Params(gt.Data(data))
	if err := crud.Create().Error(); err != nil {
		return nil, err
	}
	return data, nil
}

// update data
func (c *User) UpdateForm(params map[string][]string) interface{} {

	if err := crud.UpdateForm(params); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create data
func (c *User) CreateForm(params map[string][]string) interface{} {

	//params["createtime"] = append(params["createtime"], time.Now().Format("2006-01-02 15:04:05"))

	if err := crud.CreateForm(params); err != nil {
		//log.Log.Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
