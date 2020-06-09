package notice

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"micro-go/base-srv/models/goods/notice"
	"net/http"
)

var p notice.GsNotice

//data信息分页
//func Search(u *gin.Context) {
//	var (
//		p   record.Record
//		res interface{}
//	)
//	datas, pager, err := p.Search(cm.ToCMap(u))
//	if err != nil {
//		res = result.CError(err)
//	} else {
//		res = result.GetSuccessPager(datas, pager)
//	}
//	u.JSON(http.StatusOK, res)
//}

//新增data信息
func Create(u *gin.Context) {
	var (
		data notice.GsNotice
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
