package record

import (
	"demo/base-srv/models/record"
	cm2 "demo/commons/util/cm"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

//data信息分页
func Search(u *gin.Context) {
	var (
		p   record.Record
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
