package cron

import (
	"github.com/dreamlu/gt"
	"github.com/robfig/cron/v3"
)

// 定时任务库: https://github.com/robfig/cron
// 定时扫描删除待支付订单
func cronOrder() {
	c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc("@every 15m", deleteOrder) // 每15分钟扫描
	c.Start()
}

// 修改待付款时间大于15分钟的, 为已取消
func deleteOrder() {
	gt.Logger().Info("[定时清理订单开始]")
	err := gt.NewCrud().Select("update `order` set status = 8 where status = 0 and DATE_SUB(now(), INTERVAL 24 HOUR) > createtime").Exec().Error()
	if err != nil {
		gt.Logger().Error("[定时清理订单问题]:", err.Error())
	}
}
