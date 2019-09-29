package breaker

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	der "github.com/dreamlu/go-tool"
	"github.com/felixge/httpsnoop"
	"log"
	"net/http"
)

func init() {
	// VolumeThreshold就是单位时间内(10s)触发熔断的最低请求次数，默认为20
	// ErrorPercentThreshold即触发熔断要达到的错误率，默认为50%
	// 总结一下，即在单位时间内如果调用次数超过20次，且错误率超过50%就触发熔断
	hystrix.DefaultVolumeThreshold = 20       //0
	hystrix.DefaultErrorPercentThreshold = 50 //0
	//cl :=  micro_hystrix.NewClientWrapper()(client.DefaultClient)
}

// BreakerWrapper hystrix breaker
func BreakerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		//log.Println(name)
		err := hystrix.Do(name, func() error {
			//h.ServeHTTP(w, r)

			m := httpsnoop.CaptureMetrics(h, w, r)
			log.Println("http code:", m.Code)
			if m.Code >= http.StatusBadRequest {
				return errors.New("熔断触发")
			}

			return nil
		}, func(e error) error {
			if e == hystrix.ErrCircuitOpen {
				w.WriteHeader(http.StatusAccepted)
				//res,_ := result.GetMapData(200, "请稍后重试").String()
				_, _ = w.Write([]byte("请稍后重试"))
			}

			return e
		})
		if err != nil {
			der.Logger().Error("hystrix breaker err: ", err)
			return
		}
	})
}
