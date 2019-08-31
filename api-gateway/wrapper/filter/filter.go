package filter

import (
	"encoding/json"
	"github.com/dreamlu/go-tool/tool/result"
	"github.com/micro/micro/plugin"
	"net/http"
	"strings"
)

func Filter() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 简单验证
			if !strings.Contains(r.URL.String(), "token") {
				res, _ := json.Marshal(result.MapNoAuth)
				w.Header().Add("Content-Type", "application/json")
				_, _ = w.Write([]byte(res))
				return
			}

			// 传递下一级
			h.ServeHTTP(w, r)
			return
		})
	}
}
