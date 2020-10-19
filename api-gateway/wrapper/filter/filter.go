package filter

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/dreamlu/gt"
	"github.com/dreamlu/gt/cache"
	"github.com/dreamlu/gt/tool/result"
	"github.com/dreamlu/gt/tool/util/str"
	"github.com/dreamlu/micro/v2/plugin"
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

				// 重复点击
				switch r.Method {
				case "POST", "PATCH":
					b := check(token, path)
					if !b {
						res, _ := json.Marshal(result.GetMapData(result.CodeText, "点击太快啦"))
						w.Header().Add("Content-Type", "application/json")
						_, _ = w.Write(res)
						return
					}
				}

				h.ServeHTTP(w, r)
				return
			}

			// 传递下一级
			h.ServeHTTP(w, r)
			return
		})
	}
}

// 重复请求全局验证
func check(token, path string) bool {

	// 白名单
	if b := white(path); b {
		return true
	}
	// 判断重复下单:redis
	key := token + path
	// md5加密缩短长度key
	has := md5.Sum([]byte(key))
	key = strings.ToUpper(fmt.Sprintf("%x", has))

	ce := cache.NewCache()
	ca, _ := ce.Get(key)
	if ca.Data == nil {
		ca.Data = 1
		ca.Time = 2 * cache.CacheSecond
		_ = ce.Set(key, ca)
		//return nil
	} else {
		return false
	}
	return true
}

// 白名单
func white(path string) bool {
	switch {
	case strings.Contains(path, "/cart/create"):
		return true
	}
	if strings.Contains(path, "/upload") {
		return true
	}
	return false
}
