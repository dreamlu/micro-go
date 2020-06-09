package sms

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

// 社群特有短信发送ak
const (
	Ak           = ""
	AS           = ""
	SignName     = ""
	TemplateCode = "" // 验证码code
)

type Code struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

// 获取短信验证码
func Send(u *gin.Context) {

	var co Code
	_ = u.ShouldBindJSON(&co)

	ce := cache.NewCache()
	if cm, err := ce.Get(co.Phone); err == nil && cm.Data != nil {
		u.JSON(http.StatusOK, result.GetText("验证码5分钟内有效,无需重复发送"))
		return
	}

	code := fmt.Sprintf("%04v", rand.Int31n(10000))
	err := SendMsg(co.Phone, TemplateCode, "{\"code\":\""+code+"\"}")
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	cm := cache.CacheModel{
		Time: 5 * cache.CacheMinute,
		Data: code,
	}
	err = ce.Set(co.Phone, cm)
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	u.JSON(http.StatusOK, result.MapSuccess)
}

// 短信验证码验证
func Check(u *gin.Context) {

	var co Code
	_ = u.ShouldBindJSON(&co)

	ce := cache.NewCache()
	cm, err := ce.Get(co.Phone)
	if err == nil {
		if cm.Data.(string) == co.Code { // 验证成功
			u.JSON(http.StatusOK, result.MapSuccess)
			return
		}
		u.JSON(http.StatusOK, result.GetText("验证码不存在或已失效"))
		return
	}

	u.JSON(http.StatusOK, result.GetError(err.Error()))
}

// 发送短信
func SendMsg(phoneNumber, templateCode, TemplateParam string) error {
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", Ak, AS)
	if client == nil {
		gt.Logger().Error("短信初始化client错误")
		return errors.New("短信初始化client错误")
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumber
	request.TemplateParam = TemplateParam
	request.SignName = SignName
	request.TemplateCode = templateCode

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	fmt.Printf("response is %#v\n", response)
	return nil
}
