package goods

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-go/base-srv/models/goods"
	"micro-go/commons/util/cm"
	"net/http"
)

var p goods.Goods

//根据id获得data
func GetByID(u *gin.Context) {
	data, err := p.GetByID(cm.ToCMap(u))
	u.JSON(http.StatusOK, cm.ResGet(err, data))
}

//data信息分页
func GetBySearch(u *gin.Context) {
	ss := p.Search(cm.ToCMap(u))
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
	var data goods.GoodsM

	// json 类型需要匹配
	// 与spring boot不同
	// 不能自动将字符串转成对应类型
	// 严格匹配
	_ = u.ShouldBindJSON(&data)
	// do something

	ss := p.Update(&data)
	u.JSON(http.StatusOK, ss)
}

//新增data信息
func Create(u *gin.Context) {
	var data goods.GoodsM

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	log.Println(err)

	ss := p.Create(&data)
	u.JSON(http.StatusOK, ss)
}
