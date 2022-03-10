package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
	"time"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/china-xs/gin-tpl/middleware"
)

// Validator is a validator middleware.
func Logger(log *zap.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(c *gin.Context, req interface{}) (reply interface{}, err error) {
			var code int32
			var reason string
			startTime := time.Now()
			reply, err = handler(c, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			// 当前仅记录 到api曾的出入参数，如需独立记录额外参数，请独立配置
			log.Info("req-log",
				zap.Float64("latency", time.Since(startTime).Seconds()),
				zap.String("args", extractArgs(req)),
				zap.String("reply", extractArgs(reply)),
				zap.Int32("code", code),
				zap.String("reason", reason),
			)
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
