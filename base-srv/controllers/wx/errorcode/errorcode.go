package errorcode

// 微信通用错误
type ErrorCode struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (m ErrorCode) String() string {
	return "errcode:" + m.Errmsg + ", errmsg: " + m.Errmsg
}
