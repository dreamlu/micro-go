package delivery

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	json2 "gopkg.in/square/go-jose.v2/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// 		String key = "";				//贵司的授权key
//		String customer = "";			//贵司的查询公司编号
//		String com = "yunda";			//快递公司编码
//		String num = "3950055201640";	//快递单号
//		String phone = "";				//手机号码后四位
//		String from = "";				//出发地
//		String to = "";					//目的地
//		int resultv2 = 0;				//开启行政规划解析

const (
	key      = ""
	customer = ""
)

type DeliveryPar struct {
	Key      string `json:"key"`
	Customer string `json:"customer"`
}

func (d DeliveryPar) String() string {
	b, err := json.Marshal(d)
	if err != nil {
		gt.Logger().Error(b)
		return ""
	}
	return string(b)
}

type DeliveryCom struct {
	Com      string `json:"com"`
	Num      string `json:"num"`
	Phone    string `json:"phone"`
	From     string `json:"from"`
	To       string `json:"to"`
	Resultv2 string `json:"resultv2"`
}

func (d DeliveryCom) String() string {
	b, err := json.Marshal(d)
	if err != nil {
		gt.Logger().Error(b)
		return ""
	}
	return string(b)
}

type Delivery struct {
	DeliveryPar
	DeliveryCom
}

func (d Delivery) String() string {
	b, err := json.Marshal(d)
	if err != nil {
		gt.Logger().Error(b)
		return ""
	}
	return string(b)
}

// 快递100返回结果
//type DeliveryRes struct {
//
//}

// 快递100实时查询
func Get(u *gin.Context) {

	num := u.Query("num")
	com := u.Query("com")
	u.JSON(http.StatusOK, result.GetSuccess(Find(num, com)))
}

func Find(num, com string) (r interface{}) {
	url := "http://poll.kuaidi100.com/poll/query.do"
	var dey = &Delivery{
		DeliveryPar: DeliveryPar{
			Key:      key,
			Customer: customer,
		},
		DeliveryCom: DeliveryCom{
			Com: com,
			Num: num,
			//Phone:    "",
			//From:     "",
			//To:       "",
			Resultv2: "0",
		},
	}

	param := dey.DeliveryCom.String()
	has := md5.Sum([]byte(param + dey.Key + dey.Customer))
	sign := strings.ToUpper(fmt.Sprintf("%x", has)) //将[]byte转成16进制

	urlParam := "customer=" + dey.Customer + "&sign=" + sign + "&param=" + param
	gt.Logger().Info(urlParam)
	res, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(urlParam))
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
	_ = json2.Unmarshal(body, &r)
	return
}
