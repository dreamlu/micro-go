// @author  dreamlu
package controllers

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"micro-go/commons/util/cm"
	"micro-go/user-srv/models"
	"net/http"
)

var p models.User

//根据id获得data
func GetByID(u *gin.Context) {
	var (
		res interface{}
	)
	id := u.Query("id")
	data, err := p.GetByID(id)
	if err != nil {
		res = result.CError(err)
	}
	res = result.GetSuccess(data)
	u.JSON(http.StatusOK, res)
}

//data信息分页
func GetBySearch(u *gin.Context) {
	var (
		res interface{}
	)
	datas, pager, err := p.GetBySearch(cm.ToCMap(u))
	if err != nil {
		res = result.CError(err)
	}
	res = result.GetSuccessPager(datas, pager)
	u.JSON(http.StatusOK, res)
}

//data信息删除
func Delete(u *gin.Context) {
	var (
		res interface{}
	)
	id := u.Param("id")
	err := p.Delete(id)
	if err != nil {
		res = result.CError(err)
	}
	res = result.MapDelete
	u.JSON(http.StatusOK, res)
}

//data信息修改
func Update(u *gin.Context) {
	var (
		data models.User
		res  interface{}
	)
	// json 类型需要匹配
	// 与spring boot不同
	// 不能自动将字符串转成对应类型
	// 严格匹配
	_ = u.ShouldBindJSON(&data)
	// do something

	_, err := p.Update(&data)
	if err != nil {
		res = result.CError(err)
	}
	res = result.MapCreate
	u.JSON(http.StatusOK, res)
}

//新增data信息
func Create(u *gin.Context) {
	var (
		data models.User
		res  interface{}
	)

	// 自定义日期格式问题
	_ = u.ShouldBindJSON(&data)

	_, err := p.Create(&data)
	if err != nil {
		res = result.CError(err)
	}
	res = result.MapCreate
	u.JSON(http.StatusOK, res)
}
