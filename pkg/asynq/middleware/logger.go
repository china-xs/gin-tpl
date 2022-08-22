// Package middleware
// @author: xs
// @date: 2022/6/9
// @Description: middleware 注意需要提前设置otel 全局链路否则链路根据数据无效
package middleware

import (
	"context"
	"fmt"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"time"
)

func LoggingHandler(zapLog *zap.Logger) asynq.MiddlewareFunc {
	var trace = otel.Tracer("asynq/tasks")
	l := log.NewL(zapLog)
	return func(h asynq.Handler) asynq.Handler {
		return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
			start := time.Now()
			ctx, span := trace.Start(ctx, fmt.Sprintf("middleware-task-%v", t.Type()))
			defer span.End()
			l.With(ctx).Info("processing", zap.String("task type", t.Type()))
			err := h.ProcessTask(ctx, t)
			if err != nil {
				return err
			}
			l.With(ctx).Info("processed",
				zap.String("type", t.Type()),
				zap.String("payload", string(t.Payload())),
				zap.Error(err),
				zap.Duration("elapsed", time.Since(start)),
			)
			return nil
		})
	}
}
