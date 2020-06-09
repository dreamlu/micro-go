package wx

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/time"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp"
	"github.com/medivhzhan/weapp/payment"
	"micro-go/base-srv/controllers/wx/access_token"
	"micro-go/base-srv/models/admin/applet"
	"micro-go/base-srv/models/client"
	"micro-go/base-srv/models/order"
	"micro-go/base-srv/util/models"
	"net/http"
	time2 "time"
)

type Wx struct {
	Code string
	*applet.Applet
}

// 用户信息解密
//type WxInfo struct {
//	//*applet.Applet
//	RawData       string
//	EncryptedData string
//	Signature     string
//	IV            string
//	SessionKey    string
//}

type WxPay struct {
	TotalFee int `json:"total_fee"` // 单位分
	Openid   string
	*applet.Applet
}

type QRCode struct {
	Scene string
	Page  string
	*applet.Applet
}

// 小程序登录
func Login(u *gin.Context) {
	var (
		wx Wx
	)

	_ = u.ShouldBindJSON(&wx)
	if err := wx.GetByAppid(wx.Appid); err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}
	res, err := weapp.Login(wx.Appid, wx.Secret, wx.Code)
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	// 根据openid查找相关信息
	// 以及redis缓存等
	var cl client.Client
	gt.NewCrud(gt.Data(&cl)).
		Select(fmt.Sprintf("select %s from client where openid = ?", gt.GetColSQL(client.Client{})), res.OpenID).
		Single()
	ca := cache.NewCache()
	var model models.TokenModel
	model.ID = cl.ID
	newID, _ := id.NewID(1)
	model.Token = newID.String()
	ca.Set(model.Token, cache.CacheModel{
		Time: cache.CacheDay,
		Data: model,
	})

	u.JSON(http.StatusOK, result.GetSuccess(res).Add("id", model.ID).Add("token", model.Token).Add("admin_id", wx.AdminID))
}

// 用户信息--解密
func Info(u *gin.Context) {

	_ = u.Request.ParseForm()
	data := u.Request.Form
	//log.Println(data)
	//_, err := wx.Applet.GetByAppid(wx.Appid)
	//if err != nil {
	//	u.JSON(http.StatusOK, result.GetError(err.Error()))
	//	return
	//}
	//var info interface{}
	// 解密用户信息
	//
	// @rawData 不包括敏感信息的原始数据字符串, 用于计算签名。
	// @encryptedData 包括敏感数据在内的完整用户信息的加密数据
	// @signature 使用 sha1( rawData + session_key ) 得到字符串, 用于校验用户信息
	// @iv 加密算法的初始向量
	// @ssk 微信 session_key
	userinfo, err := weapp.DecryptUserInfo(data["raw_data"][0], data["encrypted_data"][0], data["signature"][0], data["iv"][0], data["session_key"][0])
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}
	//phone , err := weapp.DecryptPhoneNumber(session_key, encryptedData, iv)

	u.JSON(http.StatusOK, userinfo)
}

// 用户手机号--解密
func Phone(u *gin.Context) {

	_ = u.Request.ParseForm()
	data := u.Request.Form
	//log.Println(data)
	//_, err := wx.Applet.GetByAppid(wx.Appid)
	//if err != nil {
	//	u.JSON(http.StatusOK, result.GetError(err.Error()))
	//	return
	//}
	//var info interface{}
	// 解密用户信息
	//
	// @rawData 不包括敏感信息的原始数据字符串, 用于计算签名。
	// @encryptedData 包括敏感数据在内的完整用户信息的加密数据
	// @signature 使用 sha1( rawData + session_key ) 得到字符串, 用于校验用户信息
	// @iv 加密算法的初始向量
	// @ssk 微信 session_key
	//userinfo, err := weapp.DecryptUserInfo(data["raw_data"][0], data["encrypted_data"][0], data["signature"][0], data["iv"][0], data["session_key"][0])
	//if err != nil {
	//	u.JSON(http.StatusOK, result.GetError(err.Error()))
	//	return
	//}
	phone, err := weapp.DecryptPhoneNumber(data["session_key"][0], data["encrypted_data"][0], data["iv"][0])
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	u.JSON(http.StatusOK, result.GetSuccess(phone))
}

//支付,范围对应支付的多个(5)参数
func Pay(u *gin.Context) {
	var (
		wx WxPay
	)

	_ = u.ShouldBindJSON(&wx)
	if err := wx.GetByAppid(wx.Appid); err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	tradeNo, err := id.NewID(1)
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	notifyUrl := gt.Configger().GetString("app.notifyUrl") + "/pay"
	//fmt.Println("支付回调url:", notifyUrl)
	// 新建支付订单
	form := payment.Order{
		// 必填
		AppID:      wx.Appid,
		MchID:      wx.MchID,
		Body:       "商品支付",
		NotifyURL:  notifyUrl,
		OpenID:     wx.Openid,
		OutTradeNo: tradeNo.String(), //"商户订单号",
		TotalFee:   wx.TotalFee,

		// 选填 ...
		/*IP:        "发起支付终端IP",
		NoCredit:  "是否允许使用信用卡",
		StartedAt: "交易起始时间",
		ExpiredAt: "交易结束时间",
		Tag:       "订单优惠标记",
		Detail:    "商品详情",
		Attach:    "附加数据",*/
	}

	res, err := form.Unify(wx.PaySecret)
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	// 获取小程序前点调用支付接口所需参数
	params, err := payment.GetParams(res.AppID, wx.PaySecret, res.NonceStr, res.PrePayID)
	if err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}
	u.JSON(http.StatusOK, result.GetSuccess(params).Add("out_trade_no", tradeNo.String()))
}

func PayNotify(u *gin.Context) {
	fmt.Println("[回调开始]")
	//var wx applet.Applet
	//// applet参数查询
	//if err := wx.GetByAppid(wx.Appid); err != nil {
	//	u.JSON(http.StatusOK, result.GetError(err.Error()))
	//	return
	//}
	// 必须在商户平台上配置的回调地址或者发起退款时指定的 notify_url 的路由处理器下
	w := u.Writer
	req := u.Request
	// 必须在下单时指定的 notify_url 的路由处理器下
	err := payment.HandlePaidNotify(w, req, func(ntf payment.PaidNotify) (bool, string) {
		// 处理通知
		fmt.Printf("[支付结果]%#v", ntf)
		s := int8(1)
		var or = &order.Order{
			OutTradeNo: ntf.OutTradeNo,
			Paytime:    time.CTime(time2.Now()),
			Status:     &s,
		}
		err := or.UpdateNotify(or)
		if err != nil {
			fmt.Printf("[订单修改结果]%#v", err.Error())
			return false, "失败原因..." + err.Error()
		}
		//处理成功
		return true, ""
		// or
		// 处理失败 return false, "失败原因..."
	})
	if err != nil {
		gt.Logger().Error(err)
	}
	u.JSON(http.StatusOK, "回调处理完成")
}

// 获得access_token
func GetAccessToken(u *gin.Context) {

	var wx applet.Applet
	_ = u.Request.ParseForm()
	params := u.Request.Form
	if err := wx.GetByAppid(params["appid"][0]); err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}
	at := access_token.AsToken(wx.Appid, wx.Secret)
	u.JSON(http.StatusOK, result.GetSuccess(at))
}

//提现
//func WxWithDraw(u *gin.Context) {
//
//	var getinfo lib.GetInfoN
//	//生成时间戳
//	nanos := time.Now().UnixNano()
//	tradeNo := strconv.FormatInt(nanos, 10)
//
//	// 新建退款订单
//	form := payment.Transferer{
//		// 必填 ...
//		AppID:  "APPID",
//		MchID:  MchID,
//		Amount: 100, //"总金额(分)",
//		//OutRefundNo: "商户退款单号",
//		OutTradeNo: tradeNo, //"商户订单号", // or TransactionID: "微信订单号",
//		ToUser:     "ozjfE5O5hFU0cQBW4eJeaWhvIjTc",
//		Desc:       "转账描述", // 若商户传入, 会在下发给用户的退款消息中体现退款原因
//
//		/*// 选填 ...
//		IP: "发起转账端 IP 地址", // 若商户传入, 会在下发给用户的退款消息中体现退款原因
//		CheckName: "校验用户姓名选项 true/false",
//		RealName: "收款用户真实姓名", // 如果 CheckName 设置为 true 则必填用户真实姓名
//		Device:   "发起转账设备信息",*/
//	}
//
//	// 需要证书
//	res, err := form.Transfer(PaySecret, "conf/cert/apiclient_cert.pem", "conf/cert/apiclient_key.pem")
//	if err != nil {
//		u.JSON(http.StatusOK, lib.GetMapDataError(lib.CodeWxWithDraw, err.Error()))
//		return
//	}
//
//	//fmt.Printf("返回结果: %#v", res)
//	getinfo.Data = res
//	getinfo.Msg = lib.MsgSuccess
//	getinfo.Status = lib.CodeSuccess
//	u.JSON(http.StatusOK, getinfo)
//}
