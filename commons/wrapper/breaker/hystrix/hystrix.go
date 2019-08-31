package hystrix

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"log"
	"net/http"
	"strings"
)

// BreakerWrapper hystrix breaker
func BreakerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		log.Println(name)
		err := hystrix.Do(name, func() error {
			// 简单验证
			if strings.Contains(r.URL.String(), "token") {
				w.WriteHeader(http.StatusUnauthorized)
				return errors.New("熔断触发")
			}

			h.ServeHTTP(w, r)
			return nil
		}, nil)
		if err != nil {
			log.Println("hystrix breaker err: ", err)
			return
		}
	})
}