package address

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"micro-go/base-srv/models/client/address"
	"micro-go/commons/util/cm"
	"net/http"
)

var p address.ClientAddress

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
		data address.ClientAddress
	)
	// json 类型需要匹配
	// 与spring boot不同
	// 不能自动将字符串转成对应类型
	// 严格匹配
	err := u.ShouldBindJSON(&data)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}
	// do something

	_, err = p.Update(&data)
	u.JSON(http.StatusOK, cm.Res(err))
}

// 新增data信息
func Create(u *gin.Context) {
	var (
		data address.ClientAddress
	)
	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}
	_, err = p.Create(&data)
	u.JSON(http.StatusOK, cm.Res(err))
}

// 新增data信息
func CreateMore(u *gin.Context) {
	var (
		data []*address.ClientAddressP
	)
	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}
	err = p.CreateMore(data)
	u.JSON(http.StatusOK, cm.Res(err))
}
