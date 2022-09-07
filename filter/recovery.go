// Package filter
// @author: xs
// @date: 2022/9/7
// @Description: gin 中间件处理恐慌问题
package filter

import (
	"fmt"
	log2 "github.com/china-xs/gin-tpl/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

//
// Recovery gin recovery 日志记录
// @param log
// @return gin.HandlerFunc
//
func Recovery(log *zap.Logger) gin.HandlerFunc {
	l := log2.NewL(log)
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.With(c.Request.Context()).Error("recover",
					zap.String("errors", fmt.Sprintf("%v", err)),
					zap.String("debug", string(debug.Stack())),
				)
				var resp = Resp{
					Code:    500,
					Message: fmt.Sprintf("%v", err),
					Reason:  "InternalServerError",
				}
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
		}()
		c.Next()
	}
}
