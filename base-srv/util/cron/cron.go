package cron

import (
	"demo/base-srv/controllers/module/live_bro"
	"demo/base-srv/controllers/wx/access_token"
	"demo/base-srv/models/admin/applet"
	"github.com/dreamlu/gt/tool/log"
	"github.com/dreamlu/gt/tool/type/cmap"
	"time"
)

// 定时任务执行方式:
// 1.
// for{
//	time.Sleep(time.Millisecond* 100)
//	fmt.Println("Hello")
//}
// 2.
// for range time.Tick(time.Millisecond*100){
//	fmt.Println("Hello")
//}
// 3.
// c := time.Tick(5 * time.Second)
//for {
//	<- c
//	go f()
//}

// 定时任务
func Cron() {
	go cronLive()
	go cronOrder()
}

// 直播
func cronLive() {
	var (
		wx    applet.Applet
		param cmap.CMap
	)
	datas, _, err := wx.GetBySearch(param)
	if err != nil {
		log.Error(err.Error())
		return
	}
	flushLive(datas)

	// 5分钟执行一次
	for range time.Tick(time.Minute * 5) {
		go flushLive(datas)
		log.Info("[开始小程序直播列表数据刷新]")
	}
}

func flushLive(datas []*applet.Applet) {
	// 启动时执行一次
	for _, v := range datas {
		at := access_token.AsToken(v.Appid, v.Secret)
		lb := live_bro.LiveBroParam{
			AccessToken: at.AccessToken,
			LiveLimit: live_bro.LiveLimit{
				Start: 0,
				Limit: 9999,
			},
		}
		lb.Applet = v
		es, err := live_bro.FlushData(lb)
		if err != nil {
			log.Error("[err刷新直播的appid:]", v.Appid, err.Error())
			return
		}
		if es.Errcode != 0 {
			log.Error("[err刷新直播的appid:]", v.Appid, es)
			return
		}
	}
}
