package order

import (
	"demo/base-srv/models/order"
	"demo/commons/util/cm"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var p order.Order

// 批量发货excel导入
func ShipExcel(u *gin.Context) {
	file, _, err := u.Request.FormFile("file")
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}
	err = p.ShipExcel(file)
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}
	u.JSON(http.StatusOK, result.MapSuccess)
}

//data信息分页
func ExportExcel(u *gin.Context) {
	f, err := p.ExportExcel(cm.ToCMap(u))
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	u.Header("Content-Type", "application/octet-stream")
	u.Header("Content-Disposition", "attachment; filename="+"订单.xlsx")
	u.Header("Content-Transfer-Encoding", "binary")
	_ = f.Write(u.Writer)
	//u.JSON(http.StatusOK, ss)
}

//根据id获得data
func GetByID(u *gin.Context) {
	id := u.Query("id")
	data, err := p.GetByID(id)
	u.JSON(http.StatusOK, cm.ResGet(err, data))
}

//data信息分页
func GetBySearch(u *gin.Context) {
	ss := p.GetBySearch(cm.ToCMap(u))
	u.JSON(http.StatusOK, ss)
}

//data信息删除
func Delete(u *gin.Context) {
	id := u.Param("id")
	ss := p.Delete(id)
	u.JSON(http.StatusOK, ss)
}

//data信息修改
//func UpdateMore(u *gin.Context) {
//	var (
//		data []*order.Order
//		res  interface{}
//	)
//	_ = u.ShouldBindJSON(&data)
//	// do something
//
//	err := p.UpdateMore(data)
//	if err != nil {
//		res = result.CError(err)
//	} else {
//		res = result.MapUpdate
//	}
//	u.JSON(http.StatusOK, res)
//}

//data信息修改
func Update(u *gin.Context) {
	var (
		data order.Order
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
		res  interface{}
		data order.OrderGS
	)

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	log.Println(err)

	err = p.Create(&data)
	if err != nil {
		res = result.CError(err)
	} else {
		res = result.MapCreate
	}
	u.JSON(http.StatusOK, res)
}
