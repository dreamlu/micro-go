// @author  dreamlu
package routers

import (
	"github.com/dreamlu/gt/tool/util/str"
	"github.com/gin-gonic/gin"
	"micro-go/base-srv/controllers/file"
	recover2 "micro-go/commons/util/recover"
	"net/http"
	"strings"
)

var Router = SetRouter()

// router version
var V = Router.Group("/")

func SetRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	//router := gin.Default()
	router := gin.New()
	str.MaxUploadMemory = router.MaxMultipartMemory
	//router.Use(CorsMiddleware())

	// 过滤器
	router.Use(Filter())
	router.Use(recover2.Recover)
	//权限中间件
	// load the casbin model and policy from files, database is also supported.
	//e := casbin.NewEnforcer("conf/authz_model.conf", "conf/authz_policy.csv")
	//router.Use(authz.NewAuthorizer(e))

	// Ping test
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	//组的路由,version
	v1 := router.Group("/")
	{
		v := v1

		// 静态目录
		// relativePath:请求路径
		// root:静态文件所在目录
		v.Static("static", "static")
		// v.GET("/statics/file", file.StaticFile)
		//文件上传
		v.POST("/file/upload", file.UploadFile)
	}
	//不存在路由
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"msg":    "接口不存在->('.')/请求方法不存在",
		})
	})
	return router
}

// 登录失效验证
func Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ip过滤
		//if gt.Configger().GetString("app.devMode") == str.Dev {
		//	c.Next()
		//	return
		//}
		//if c.Request.Method == "GET" {
		//	c.Next()
		//	return
		//}
		path := c.Request.URL.String()

		// 静态服务器 file 处理
		if strings.Contains(path, "/static/file/") {
			file.StaticFile(c)
			c.Abort()
			return
		}

		//if !strings.Contains(path, "login") &&
		//	!strings.Contains(path, "/static/file") &&
		//	!strings.Contains(path, "/notify") {
		//	token := c.GetHeader("token")
		//	if token == "" {
		//		c.Abort()
		//		c.JSON(http.StatusOK, result.GetMapData(result.CodeError, "缺少token"))
		//		return
		//	}
		//	ca := cache.NewCache()
		//	log.Println("[token]:", token)
		//	cam, err := ca.Get(token)
		//
		//	if err != nil {
		//		c.Abort()
		//		c.JSON(http.StatusOK, result.GetMapData(result.CodeError, "token不合法"))
		//		return
		//	}
		//	// 延长token对应时间
		//	ca.Set(token, cam)
		//	c.Next()
		//}
	}
}

// 处理跨域请求,支持options访问
//func Cors() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		method := c.Request.Method
//		//fmt.Println(method)
//		c.Header("Access-Control-Allow-Origin", "*")
//		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
//		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PUT, DELETE")
//		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
//		c.Header("Access-Control-Allow-Credentials", "true")
//
//		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
//		if method == "OPTIONS" {
//			c.AbortWithStatus(http.StatusNoContent)
//		}
//		// 处理请求
//		c.Next()
//	}
//}
