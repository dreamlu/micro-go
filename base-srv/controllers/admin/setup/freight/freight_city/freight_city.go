package freight_city

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"micro-go/base-srv/models/admin/setup/freight/freight_city"
	cm2 "micro-go/commons/util/cm"
	"net/http"
)

var p freight_city.FreightCity

//根据id获得data
func Get(u *gin.Context) {
	var (
		res interface{}
	)
	id := u.Query("id")
	data, err := p.Get(id)
	if err != nil {
		res = result.CError(err)
	} else {
		res = result.GetSuccess(data)
	}
	u.JSON(http.StatusOK, res)
}

//data信息分页
func Search(u *gin.Context) {
	var (
		res interface{}
	)
	datas, pager, err := p.Search(cm2.ToCMap(u))
	if err != nil {
		res = result.CError(err)
	} else {
		res = result.GetSuccessPager(datas, pager)
	}
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
	} else {
		res = result.MapDelete
	}
	u.JSON(http.StatusOK, res)
}

//data信息修改
func Update(u *gin.Context) {
	var (
		data freight_city.FreightCity
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
	} else {
		res = result.MapUpdate
	}
	u.JSON(http.StatusOK, res)
}

//新增data信息
func Create(u *gin.Context) {
	var (
		data freight_city.FreightCity
		res  interface{}
	)

	// 自定义日期格式问题
	_ = u.ShouldBindJSON(&data)

	_, err := p.Create(&data)
	if err != nil {
		res = result.CError(err)
	} else {
		res = result.MapCreate
	}
	u.JSON(http.StatusOK, res)
}

func CreateAll(u *gin.Context) {
	var (
		data []*freight_city.FreightCityPar
		res  interface{}
	)

	// 自定义日期格式问题
	_ = u.ShouldBindJSON(&data)

	_, err := p.CreateAll(data)
	if err != nil {
		res = result.CError(err)
	} else {
		res = result.MapCreate
	}
	u.JSON(http.StatusOK, res)
}
