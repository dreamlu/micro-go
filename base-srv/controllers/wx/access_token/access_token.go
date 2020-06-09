package access_token

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//小程序access_token
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// 获得access_token
func AsToken(appid, secret string) AccessToken {
	te_uri := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret
	res, _ := http.Get(te_uri)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Printf("返回结果: %#v", res)
	var at AccessToken
	_ = json.Unmarshal(body, &at)
	return at
}
