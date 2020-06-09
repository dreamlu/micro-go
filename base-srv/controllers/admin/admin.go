package admin

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/util"
	"github.com/gin-gonic/gin"
	"log"
	"micro-go/base-srv/models/admin"
	"micro-go/base-srv/util/models"
	cm2 "micro-go/commons/util/cm"
	"net/http"
)

var p admin.Admin

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
	var data admin.Admin

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
	var data admin.Admin

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	log.Println(err)

	ss := p.Create(&data)
	u.JSON(http.StatusOK, ss)
}

// 登录
func Login(u *gin.Context) {
	var login, sqlData admin.Admin

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&login)
	log.Println(err)

	// 查找
	gt.NewCrud(gt.Data(&sqlData)).Select("select id,password,role from admin where name = ?", login.Name).Single()

	// 验证不通过
	if sqlData.Password != util.AesEn(login.Password) {
		u.JSON(http.StatusOK, result.MapCountErr)
		return
	}

	ca := cache.NewCache()
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}
	var model models.TokenModel
	model.ID = sqlData.ID
	newID, _ := id.NewID(1)
	model.Token = newID.String()
	ca.Set(model.Token, cache.CacheModel{
		Time: cache.CacheDay,
		Data: model,
	})

	u.JSON(http.StatusOK, result.MapSuccess.
		Add("id", model.ID).
		Add("token", model.Token).
		Add("role", sqlData.Role))
}
