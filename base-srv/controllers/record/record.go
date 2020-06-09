package record

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"micro-go/base-srv/models/record"
	cm2 "micro-go/commons/util/cm"
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
