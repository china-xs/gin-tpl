package logger

import (
	"fmt"
	plog "github.com/china-xs/gin-tpl/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
	"time"

	"github.com/china-xs/gin-tpl/middleware"
)

// Logger Validator is a validator middleware.
func Logger(log *zap.Logger) middleware.Middleware {
	l := plog.NewL(log)
	return func(handler middleware.Handler) middleware.Handler {
		return func(c *gin.Context, req interface{}) (reply interface{}, err error) {
			var code int32
			var reason string
			var body string
			startTime := time.Now()
			reply, err = handler(c, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			// 记录body 数据
			if bodyBytes, ok := c.Get(gin.BodyBytesKey); ok {
				body = string(bodyBytes.([]byte))
			}

			var fields = make([]zap.Field, 10)
			fields[0] = zap.String("url", c.Request.URL.String())
			fields[1] = zap.String("method", c.Request.Method)
			fields[2] = zap.String("body", body)
			fields[3] = zap.String("host", c.Request.Host)
			fields[4] = zap.String("ipv4", c.ClientIP())
			fields[5] = zap.String("latency", time.Since(startTime).String())
			fields[6] = zap.String("args", extractArgs(req))
			fields[7] = zap.String("reply", extractArgs(reply))
			fields[8] = zap.Int32("code", code)
			fields[9] = zap.String("reason", reason)
			// 当前仅记录 到api曾的出入参数，如需独立记录额外参数，请独立配置
			l.With(c.Request.Context()).Info("req-log", fields...)
			return
		}
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
