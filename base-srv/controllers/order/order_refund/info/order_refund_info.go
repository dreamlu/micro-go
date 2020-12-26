package info

import (
	"demo/base-srv/models/order/order_refund/info"
	"demo/commons/util/cm"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

var p info.OrderRInfo

//根据id获得data
func Get(u *gin.Context) {
	data, err := p.Get(cm.ToCMap(u))
	u.JSON(http.StatusOK, cm.ResGet(err, data))
}

//data信息分页
func Search(u *gin.Context) {
	datas, pager, err := p.Search(cm.ToCMap(u))
	u.JSON(http.StatusOK, cm.ResPager(err, datas, pager))
}

//data信息删除
func Delete(u *gin.Context) {
	id := u.Param("id")
	err := p.Delete(id)
	u.JSON(http.StatusOK, cm.Res(err))
}

//data信息修改
func Update(u *gin.Context) {
	var (
		data info.OrderRInfo
	)
	err := u.ShouldBindJSON(&data)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}
	// do something

	_, err = p.Update(&data)
	u.JSON(http.StatusOK, cm.Res(err))
}

//新增data信息
func Create(u *gin.Context) {
	var (
		data info.OrderRInfo
	)
	err := u.ShouldBindJSON(&data)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}
	_, err = p.Create(&data)
	u.JSON(http.StatusOK, cm.Res(err))
}
