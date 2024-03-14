package middleware

import (
	log "coderunner/logger"
	"coderunner/model"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorRecover 异常恢复中间件
func ErrorRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				// 自定义类型
				case model.RespType:
					log.Logger.Warnf(
						"Request Fail by recover: url=[%s], resp=[%+v]", c.Request.URL.Path, v)
					var data interface{}
					if v.Data() == nil {
						data = []string{}
					}
					model.Result(c, v, data)
				// 其他类型
				default:
					log.Logger.Errorf("stacktrace from panic: %+v\n%s", r, string(debug.Stack()))
					model.Fail(c, model.SystemError)
				}
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
