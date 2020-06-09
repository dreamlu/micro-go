package wx

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/payment"
	"log"
	"micro-go/base-srv/models/admin/applet"
	"micro-go/base-srv/models/order"
	"micro-go/base-srv/models/order/order_refund"
	"net/http"
	"strconv"
	"strings"
)

type WxOrder struct {
	AdminID       uint64 `json:"admin_id"`
	OutTradeNo    string `json:"out_trade_no"`    // 支付单号
	OrderRefundID uint64 `json:"order_refund_id"` // 订单退款id
	applet.Applet
}

// 退款
func Refund(u *gin.Context) {
	var (
		wx       WxOrder
		or       order.Order
		of       order_refund.OrderRefund
		allMoney float64
	)
	_ = u.ShouldBindJSON(&wx)

	// order参数查询
	datas, err := or.GetByOutTradeNo(wx.OutTradeNo)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}
	for _, v := range datas {
		allMoney += v.Money
	}

	// 退款金额查询
	of, err = of.GetByID(wx.OrderRefundID)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}

	// applet参数查询
	if err = wx.GetByAdminID(wx.AdminID); err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}

	refundNo := of.OutRefundNo + "-" + strconv.FormatUint(wx.OrderRefundID, 10)
	notifyUrl := gt.Configger().GetString("app.notifyUrl") + "/refund/" + wx.PaySecret
	// 新建退款订单
	form := payment.Refunder{
		// 必填
		AppID:       wx.Appid,
		MchID:       wx.MchID,
		TotalFee:    int(allMoney * 100), //"总金额(分)"
		RefundFee:   int(of.Money * 100), //"退款金额(分)"
		OutRefundNo: refundNo,
		// 二选一
		OutTradeNo: wx.OutTradeNo, // or TransactionID: "微信订单号",
		// 选填 ...
		RefundDesc: "用户退款",    // 若商户传入, 会在下发给用户的退款消息中体现退款原因
		NotifyURL:  notifyUrl, //结果通知地址，覆盖商户平台上配置的回调地址
	}

	// 需要证书
	res, err := form.Refund(wx.PaySecret, wx.AppCert, wx.AppKey)
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}
	log.Printf("返回结果: %#v", res)
	//err = o.EditStatus(id, 4)
	//if err != nil {
	//	u.JSON(http.StatusOK, result.GetMapData(500, "退款失败"))
	//	return
	//}
	u.JSON(http.StatusOK, result.MapSuccess)
}

func RefundNotify(u *gin.Context) {
	fmt.Println("退款回调开始")
	// 简单点,直接用secret
	secret := u.Param("secret")
	gt.Logger().Info("secret:", secret)
	//var wx applet.Applet
	// applet参数查询
	//if err := wx.GetByAdminID(admin_id); err != nil {
	//	u.JSON(http.StatusOK, result.GetError(err.Error()))
	//	return
	//}
	// 必须在商户平台上配置的回调地址或者发起退款时指定的 notify_url 的路由处理器下
	w := u.Writer
	req := u.Request
	err := payment.HandleRefundedNotify(w, req, secret, func(notify payment.RefundedNotify) (b bool, s string) {
		status := notify.RefundStatus
		fmt.Println("[退款回调状态]:", status, " 退款单号", notify.OutRefundNo)
		if status == "SUCCESS" {
			var of = order_refund.OrderRefund{
				Status: 4,
			}
			i := strings.Split(notify.OutRefundNo, "-")[1]
			of.ID, _ = strconv.ParseUint(i, 10, 64)

			err := of.Update(&of)
			if err != nil {
				return false, "失败原因..." + err.Error()
			}
			return true, ""
		}
		var msg string
		if status == "CHANGE" {
			msg = "退款异常"
		} else {
			msg = status
		}
		gt.Logger().Error(msg)
		return false, msg
	})
	if err != nil {
		gt.Logger().Error(err)
	}
}
