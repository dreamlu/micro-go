package live_bro

import (
	"bytes"
	"encoding/json"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"micro-go/base-srv/models/admin/applet"
	"micro-go/base-srv/models/module/live_bro"
	"net/http"
)

// 小程序直播
// 开发文档:https://docs.qq.com/doc/DZHhzV0FiYXRQV01i
// 使用文档:https://res.wx.qq.com/mmbizwxampnodelogicsvr_node/dist/images/help_0f7865.pdf
type LiveBroParam struct {
	*applet.Applet
	AccessToken string `json:"access_token"`
	LiveLimit
}

type LiveLimit struct {
	// "start": 0, // 起始拉取房间，start=0表示从第1个房间开始拉取
	// "limit": 10 // 每次拉取的个数上限，不要设置过大，建议100以内
	Start int64 `json:"start"`
	Limit int64 `json:"limit"`
}

type ErrorRes struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//根据id获得data
func LiveBroList(u *gin.Context) {
	var (
		param LiveBroParam
		res   interface{}
		lb    live_bro.LiveBroData
	)
	_ = u.ShouldBindJSON(&param)
	// 默认从数据库查找数据列表
	data, err := lb.Get(param.AdminID)
	if data == nil || data.ID == 0 {
		es, err := FlushData(param)
		if err != nil {
			u.JSON(http.StatusOK, result.CError(err))
			return
		}
		if es.Errcode != 0 {
			u.JSON(http.StatusOK, es)
			return
		}
	}
	if err != nil {
		res = result.CError(err)
	} else {
		res = result.GetSuccess(data)
	}
	u.JSON(http.StatusOK, res)
}

func FlushData(param LiveBroParam) (es *ErrorRes, err error) {
	var (
		lb  live_bro.LiveBroData
		uri = "http://api.weixin.qq.com/wxa/business/getliveinfo?access_token="
	)
	p, err := json.Marshal(param.LiveLimit)
	if err != nil {
		return nil, err
	}
	res, _ := http.Post(uri+param.AccessToken, "application/json", bytes.NewReader(p))
	resData, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(resData, &es)
	if es.Errcode != 0 {
		return es, nil
	}
	//log.Println("[直播线上拉取结果]", res)
	// 将数据存入数据库
	lb.Data = resData
	// 默认从数据库查找数据列表
	data, _ := lb.Get(param.AdminID)
	if data != nil && data.ID != 0 { // 已存在数据,进行更新
		_, err = lb.Update(&lb)
		return
	}

	lb.AdminID = param.AdminID
	_, err = lb.Create(&lb)
	if err != nil {
		return nil, err
	}
	return
}
