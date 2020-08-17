package recover

import (
	"github.com/dreamlu/gt/tool/result"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"strings"
)

// 异常捕获
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			//log.Printf("panic: %v\n", r)
			debug.PrintStack()
			//封装通用json返回
			//c.JSON(http.StatusOK, Result.Fail(errorToString(r)))
			//Result.Fail不是本例的重点，因此用下面代码代替
			ss := strings.Split(string(debug.Stack()), "\n\t")
			res := make(map[string]string)
			for _, v := range ss {
				ks := strings.Split(v, "\n")
				res[ks[0]] = ks[1]
			}
			c.JSON(http.StatusOK, result.GetError(res))
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
