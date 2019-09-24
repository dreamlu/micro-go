// @author  dreamlu
package models

import (
	der "github.com/dreamlu/go-tool"
	"github.com/dreamlu/go-tool/tool/result"
	time2 "github.com/dreamlu/go-tool/tool/type/time"
	"time"
)

/*user model*/
type User struct {
	ID         uint        `json:"id" gorm:"primary_key"`
	Name       string      `json:"name" valid:"required,len=2-20"`
	Createtime time2.CTime `json:"createtime"` //maybe you like CDate
}

var crud = der.NewCrud(
	&der.CrudParam{
		Model: User{}, // model
		Table: "user", // table name
	})

// get user, by id
func (c *User) GetByID(id string) interface{} {

	var user User // not use *User
	crud.Param().ModelData = &user
	if err := crud.GetByID(id); err != nil {
		der.Logger().Error(err.Error())
		return result.GetError(err.Error())
	}
	der.Logger().Info("测试", crud.Param())
	return result.GetSuccess(user)
}

// get user, limit and search
// clientPage 1, everyPage 10 default
func (c *User) GetBySearch(params map[string][]string) interface{} {
	var users []*User
	crud.Param().ModelData = &users

	pager, err := crud.GetBySearch(params)
	if err != nil {
		der.Logger().Error(err.Error())
		return result.GetError(err)
	}
	return result.GetSuccessPager(users, pager)
}

// delete user, by id
func (c *User) Delete(id string) interface{} {

	if err := crud.Delete(id); err != nil {
		der.Logger().Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeDelete, result.MsgDelete)
}

// update user
func (c *User) Update(data *User) interface{} {

	if err := crud.Update(data); err != nil {
		der.Logger().Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeUpdate, result.MsgUpdate)
}

// create user
func (c *User) Create(data *User) interface{} {

	// create time
	(*data).Createtime = time2.CTime(time.Now())

	if err := crud.Create(data); err != nil {
		der.Logger().Error(err.Error())
		return result.GetError(err)
	}
	return result.GetMapData(result.CodeCreate, result.MsgCreate)
}
