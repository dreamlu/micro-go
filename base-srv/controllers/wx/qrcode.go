package wx

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/id"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp/code"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"micro-go/base-srv/controllers/wx/access_token"
	"micro-go/base-srv/models/admin/applet"
	"micro-go/commons/models"
	"micro-go/commons/util/cm"
	"net/http"
)

// TODO 约定个时间久的删除key-value, 定时任务, 同时清理一次性key
/// ========= 二维码 ===============
type QrCode struct {
	models.IDCom
	Key   string `json:"key" gorm:"varchar(20)"`           // key
	Value string `json:"value"`                            // value
	Type  byte   `json:"type" gorm:"tinyint(1);DEFAULT:0"` // 0一直存在,1一次性key
}

// 创建二维码key==value转换值
func CreateQrCode(qc *QrCode) {
	//var qc QrCode
	//_ = u.ShouldBindJSON(&qc)
	if qc.Value != "" {
		ids, _ := id.NewID(1)
		qc.Key = ids.String()
	}
	cd := gt.NewCrud(
		gt.Model(QrCode{}),
		gt.Data(&qc),
	).Create()
	if cd.Error() != nil {
		return
	}
}

// 查找
func (c *QrCode) getByValue() (*QrCode, error) {
	cd := gt.NewCrud(
		gt.Data(&c),
	).Select("select `key` from qr_code where value = ?", c.Value).Single()

	if cd.Error() != nil {
		return nil, cd.Error()
	}

	return c, nil
}

func GetByKey(u *gin.Context) {
	var (
		params = cmap.CMap{}
		qc     QrCode
	)

	params = cm.ToCMap(u)
	cd := gt.NewCrud(
		gt.Model(QrCode{}),
		gt.Data(&qc),
	).GetByData(params)
	if cd.Error() != nil {
		u.JSON(http.StatusOK, result.CError(cd.Error()))
		return
	}

	// 请求成功后, 如果是一次性key, 则删除对应状态
	if qc.ID != 0 && qc.Type == 1 {
		err := cd.Select("delete from `qr_code` where id = ?", qc.ID).Exec().Error()
		if err != nil {
			gt.Logger().Error("二维码key删除问题", err.Error())
		}
	}

	u.JSON(http.StatusOK, result.GetSuccess(qc.Value))
}

//二维码,业务量多的情况
func GetQRCode(u *gin.Context) {
	var wx applet.Applet
	_ = u.Request.ParseForm()
	params := u.Request.Form
	if err := wx.GetByAppid(params["appid"][0]); err != nil {
		u.JSON(http.StatusOK, result.GetError(err.Error()))
		return
	}

	// scene参数转换
	qc := &QrCode{}
	qc.Value = params["value"][0]
	if p, ok := params["type"]; ok && p[0] == "1" {
		qc.Type = 1 // 一次性key
	}
	_, err := qc.getByValue()
	if err != nil {
		if err.Error() == result.MsgNoResult {
			CreateQrCode(qc)
			goto into
		}
		u.JSON(http.StatusOK, result.CError(err))
		return
	}

into:
	coder := code.QRCoder{
		Scene:     qc.Key,            // 参数数据
		Page:      params["page"][0], // 识别二维码后进入小程序的页面链接
		Width:     430,               // 图片宽度
		IsHyaline: false,             // 是否需要透明底色
		AutoColor: true,              // 自动配置线条颜色, 如果颜色依然是黑色, 则说明不建议配置主色调
		LineColor: code.Color{ //  AutoColor 为 false 时生效, 使用 rgb 设置颜色 十进制表示
			R: "50",
			G: "50",
			B: "50",
		},
	}

	at := access_token.AsToken(wx.Appid, wx.Secret)
	// token: 微信 access_token
	resu, err := coder.UnlimitedAppCode(at.AccessToken)
	defer resu.Body.Close()
	if err != nil {
		u.JSON(http.StatusOK, result.CError(err))
		return
	}

	bodyu, _ := ioutil.ReadAll(resu.Body)
	u.Writer.Header().Add("Content-Type", "image/png")
	u.Writer.Write(bodyu)
}

// 普通二维码生成
func PQrcode(u *gin.Context) {
	url := u.Query("url")
	w := u.Writer

	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}()
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(png)))
	w.Write(png)
}
