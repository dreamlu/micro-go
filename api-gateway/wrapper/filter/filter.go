package filter

import (
	"encoding/json"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/util/str"
	"github.com/micro/micro/v2/plugin"
	"log"
	"net/http"
	"strings"
)

// 全局拦截器
func Filter() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// ip过滤
			if gt.Configger().GetString("app.devMode") == str.Dev &&
				strings.Contains(r.RemoteAddr, "127.0.0.1") {
				// 传递下一级
				h.ServeHTTP(w, r)
				return
			}

			if r.Method == "GET" {
				h.ServeHTTP(w, r)
				return
			}

			path := r.URL.String()
			//fmt.Println("请求url:",path)
			if !strings.Contains(path, "login") &&
				!strings.Contains(path, "/static/file") &&
				!strings.Contains(path, "/notify") &&
				!strings.Contains(path, "/live_bro") {
				token := r.Header.Get("token")
				if token == "" {
					res, _ := json.Marshal(result.GetMapData(result.CodeError, "缺少token"))
					w.Header().Add("Content-Type", "application/json")
					_, _ = w.Write(res)
					return
				}
				ca := cache.NewCache()
				log.Println("[token]:", token)
				cam, err := ca.Get(token)

				if err != nil {
					res, _ := json.Marshal(result.GetMapData(result.CodeError, "token不合法"))
					w.Header().Add("Content-Type", "application/json")
					_, _ = w.Write(res)
					return
				}
				// 延长token对应时间
				_ = ca.Set(token, cam)
				h.ServeHTTP(w, r)
				return
			}

			// 传递下一级
			h.ServeHTTP(w, r)
			return
		})
	}
}
