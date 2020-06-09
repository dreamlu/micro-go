package template_msg

import (
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/type/json"
	json2 "gopkg.in/square/go-jose.v2/json"
	"io/ioutil"
	"micro-go/base-srv/controllers/wx/access_token"
	"micro-go/base-srv/controllers/wx/errorcode"
	"micro-go/base-srv/models/admin/applet"
	"net/http"
	"strings"
)

// 模板消息
type ModelMsg struct {
	//AccessToken string     `json:"access_token"`
	Touser     string     `json:"touser" gorm:"type:varchar(50)"`      // 接收者（用户）的 openid
	TemplateID string     `json:"template_id" gorm:"type:varchar(50)"` // 模板id
	Page       string     `json:"page" gorm:"type:varchar(50)"`        // 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	Data       json.CJSON `json:"data" gorm:"type:json"`
	// miniprogram_state 跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
}

func (m ModelMsg) String() string {
	b, err := json2.Marshal(m)
	if err != nil {
		gt.Logger().Error(b)
		return ""
	}
	return string(b)
}

// 小程序订阅消息-发送
// 目前小程序只支持一条条发送
func (m ModelMsg) Send(adminID uint64) {

	var wx applet.Applet
	if err := wx.GetByAdminID(adminID); err != nil {
		return
	}
	at := access_token.AsToken(wx.Appid, wx.Secret)

	url := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + at.AccessToken
	res, err := http.Post(url, "application/json", strings.NewReader(m.String()))
	if err != nil {
		gt.Logger().Error(err.Error())
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		gt.Logger().Error(err.Error())
		return
	}
	var er errorcode.ErrorCode
	_ = json2.Unmarshal(body, &er)
	if er.Errcode != 0 {
		gt.Logger().Error("订阅消息发送失败: ", er.String())
		return
	}
	return
}
