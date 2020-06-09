package pca

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func Ws(u *gin.Context) {
	key := u.Query("key")
	output := u.Query("output")
	url := "https://apis.map.qq.com/ws/district/v1/list?output=" + output + "&key=" + key
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)
	var data interface{}
	_ = json.Unmarshal(body, &data)
	u.JSON(http.StatusOK, data)
}
