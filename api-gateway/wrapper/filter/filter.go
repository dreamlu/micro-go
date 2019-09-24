package filter

import (
	"github.com/micro/micro/plugin"
	"net/http"
)

// 全局拦截器
func Filter() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 简单验证
			// 写自己的全局验证逻辑
			//if !strings.Contains(r.URL.String(), "token") {
			//	res, _ := json.Marshal(result.MapNoAuth)
			//	w.Header().Add("Content-Type", "application/json")
			//	_, _ = w.Write(res)
			//	return
			//}

			// 传递下一级
			h.ServeHTTP(w, r)
			return
		})
	}
}
