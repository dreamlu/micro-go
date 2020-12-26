package order_refund

import (
	"demo/base-srv/models/order/order_refund"
	cm2 "demo/commons/util/cm"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var p order_refund.OrderRefund

//根据id获得data
func GetByID(u *gin.Context) {
	id := u.Query("id")
	data, err := p.GetByID(id)
	u.JSON(http.StatusOK, cm2.ResGet(err, data))
}

//data信息分页
func GetBySearch(u *gin.Context) {
	ss := p.GetBySearch(cm2.ToCMap(u))
	u.JSON(http.StatusOK, ss)
}

//data信息删除
func Delete(u *gin.Context) {
	id := u.Param("id")
	ss := p.Delete(id)
	u.JSON(http.StatusOK, ss)
}

//data信息修改
func Update(u *gin.Context) {
	var (
		res  interface{}
		data order_refund.OrderRefund
	)

	// json 类型需要匹配
	// 与spring boot不同
	// 不能自动将字符串转成对应类型
	// 严格匹配
	_ = u.ShouldBindJSON(&data)
	// do something

	err := p.Update(&data)
	if err != nil {
		res = result.CError(err)
	} else {
		res = result.MapUpdate
	}
	u.JSON(http.StatusOK, res)
}

//新增data信息
func Create(u *gin.Context) {
	var data order_refund.OrderRefund

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	log.Println(err)

	ss := p.Create(&data)
	u.JSON(http.StatusOK, ss)
}
