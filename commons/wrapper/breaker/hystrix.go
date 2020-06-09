package breaker

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	http2 "micro-go/commons/wrapper/http"
	"net/http"
)

func init() {
	// VolumeThreshold就是单位时间内(10s)触发熔断的最低请求次数，默认为20
	// ErrorPercentThreshold即触发熔断要达到的错误率，默认为50%
	// 总结一下，即在单位时间内如果调用次数超过20次，且错误率超过50%就触发熔断
	hystrix.DefaultVolumeThreshold = 1000     //0
	hystrix.DefaultErrorPercentThreshold = 50 //0
	//cl :=  micro_hystrix.NewClientWrapper()(client.DefaultClient)
}

// BreakerWrapper hystrix breaker
func BreakerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		hystrix.Do(name, func() error {
			sct := &http2.StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}
			h.ServeHTTP(sct.WrappedResponseWriter(), r)

			if sct.Status >= http.StatusInternalServerError {
				str := fmt.Sprintf("status code %d", sct.Status)
				return errors.New(str)
			}
			return nil
		}, func(e error) error {
			if e == hystrix.ErrCircuitOpen {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("请稍后重试"))
			}
			return e
		})
	})
}
