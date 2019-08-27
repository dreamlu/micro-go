// @author  dreamlu
package controllers

import (
	"github.com/dreamlu/go-tool/tool/xss"
	"github.com/gin-gonic/gin"
	"log"
	"micro-go/user-srv/models"
	"net/http"
)

var p models.User

// 根据id获得用户获取
func GetById(u *gin.Context) {
	id := u.Query("id")
	ss := p.GetByID(id)
	u.JSON(http.StatusOK, ss)
}

// 用户信息分页
func GetBySearch(u *gin.Context) {
	// this is get url 参数
	_ = u.Request.ParseForm()
	values := u.Request.Form //在使用之前需要调用ParseForm方法
	xss.XssMap(values)
	ss := p.GetBySearch(values)
	u.JSON(http.StatusOK, ss)
}

// 用户信息删除
func Delete(u *gin.Context) {
	id := u.Param("id")
	ss := p.Delete(id)
	u.JSON(http.StatusOK, ss)
}

// 用户信息修改
func Update(u *gin.Context) {
	var user models.User

	// json 类型需要匹配
	// 与spring boot不同
	// 不能自动将字符串转成对应类型
	// 严格匹配
	_ = u.ShouldBindJSON(&user)
	// do something

	ss := p.Update(&user)
	u.JSON(http.StatusOK, ss)
}

// 新增用户信息
func Create(u *gin.Context) {
	var user models.User

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&user)
	log.Println(err)

	ss := p.Create(&user)
	u.JSON(http.StatusOK, ss)
}
