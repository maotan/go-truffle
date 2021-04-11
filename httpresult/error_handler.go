/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package httpresult

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			log.Infof("panic: %#v", r)
			//debug.PrintStack()
			switch v := r.(type) {
			case WarnError:
				c.JSON(http.StatusOK, gin.H{
					"code": v.Code,
					"msg":  v.Msg,
					"data": nil,
				})
			case error:
				c.JSON(http.StatusOK, Fail(-1, v.Error()))
				debug.PrintStack()
			default:
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  r.(string),
					"data": nil,
				})
				debug.PrintStack()
			}

			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
