package cron

import (
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/tool/type/cmap"
	"github.com/dreamlu/gt/tool/type/time"
	"github.com/robfig/cron/v3"
	"micro-go/base-srv/models/goods/notice"
	time2 "time"
)

func CTimeCron(t time.CTime) (rs string) {
	// "2006 01 02 15 04 05"
	s := time2.Time(t).Format("05 04 15 02 01 *")
	//sts := strings.Split(s, " ")
	//for _, v := range sts {
	//	if strings.HasPrefix(v, "0") {
	//		v = string([]byte(v)[1:])
	//	}
	//	rs += v + " "
	//}
	//rs = string([]byte(rs)[:len(rs)-1])
	return s
}

// 开团通知
// 根据传入的时间创建定时任务
func CronNotice(t time.CTime) {
	c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc(CTimeCron(t), NoticeInfo)
	gt.Logger().Info(CTimeCron(t))
	c.Start()
}

func NoticeInfo() {
	fmt.Println("[定时任务:开团通知查询开始]")

	// 查找开团预定中,未通知的
	// 小于等于当前时间的商品开始时间
	var (
		params = cmap.CMap{}
		nts    []*notice.GsNotice
	)
	params.Add("status", "0")
	whereSQL := "now() >= start_time"
	cd := gt.NewCrud(
		gt.Model(notice.GsNotice{}),
		gt.Data(&nts),
		gt.InnerTable([]string{"gs_notice", "goods"}),
		gt.SubWhereSQL(whereSQL),
	).GetMoreBySearch(params)
	if cd.Error() != nil {
		gt.Logger().Error(cd.Error().Error())
		return
	}

	// 1.进行通知
	// 2.通知后修改对应通知状态
	var goodsIDS []uint64
	for _, v := range nts {
		v.Send(v.AdminID)
		for _, v2 := range goodsIDS {
			if v2 == v.GoodsID {
				goto into
			}
		}
		goodsIDS = append(goodsIDS, v.GoodsID)
	into:
	}

	// 修改状态
	gt.Logger().Info("[商品id]:", goodsIDS)
	for _, v := range goodsIDS {
		err := gt.NewCrud().Select("update gs_notice set status = 1 where goods_id = ?", v).Exec()
		if err != nil {
			gt.Logger().Error(err.Error())
		}
	}

}
