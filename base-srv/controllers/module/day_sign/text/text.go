package text

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-go/base-srv/models/module/day_sign/text"
	cm2 "micro-go/commons/util/cm"
	"net/http"
)

var p text.DaySignText

//根据id获得data
func GetByID(u *gin.Context) {
	id := u.Query("id")
	ss := p.GetByID(id)
	u.JSON(http.StatusOK, ss)
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
	var data text.DaySignText

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
	var data text.DaySignText

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	log.Println(err)

	ss := p.Create(&data)
	u.JSON(http.StatusOK, ss)
}
