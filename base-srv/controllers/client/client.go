package client

import (
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"log"
	"micro-go/base-srv/models/client"
	"micro-go/base-srv/util/models"
	"micro-go/commons/util/cm"
	"net/http"
	"strconv"
)

var p client.Client

// token
func Token(u *gin.Context) {
	client_id := u.Query("id")
	if client_id == "" {
		u.JSON(http.StatusOK, result.GetError("id不能为空"))
		return
	}
	ca := cache.NewCache()
	var model models.TokenModel
	model.ID, _ = strconv.ParseUint(client_id, 10, 64)
	newID, _ := id.NewID(1)
	model.Token = newID.String()
	ca.Set(model.Token, cache.CacheModel{
		Time: cache.CacheDay,
		Data: model,
	})

	u.JSON(http.StatusOK, result.MapSuccess.Add("token", model.Token))
}

//根据id获得data
func GetByIDC(u *gin.Context) {
	data, err := p.GetByIDC(cm.ToCMap(u))
	u.JSON(http.StatusOK, cm.ResGet(err, data))
}

//根据id获得data
func GetByID(u *gin.Context) {
	data, err := p.GetByID(cm.ToCMap(u))
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
func Update(u *gin.Context) {
	var data client.Client

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
	var data client.Client

	// 自定义日期格式问题
	err := u.ShouldBindJSON(&data)
	log.Println(err)

	ss := p.Create(&data)
	u.JSON(http.StatusOK, ss)
}
