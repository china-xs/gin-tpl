// Package redis hook钩子
// @author: xs
// @date: 2022/8/3
// @Description: redis
package redis

import (
	"bytes"
	"context"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type Hook struct {
	//log *zap.Logger
	log *log.Log
}

//
// NewHook 初始化redis 钩子
// @param log
//
func NewHook(zl *zap.Logger) redis.Hook {
	l := log.NewL(zl)
	return &Hook{
		log: l,
	}
}

const ctxKey = "ctxKey"

//
// BeforeProcess 执行前
// @receiver h
// @param ctx
// @param cmd
// @return context.Context
// @return error
//
func (h Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	t := time.Now()
	ctx = context.WithValue(ctx, ctxKey, t)
	return ctx, nil
}

//
// AfterProcess 单个命令执行后
// @receiver h
// @param ctx
// @param cmd
// @return error
//
func (h Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	v := ctx.Value(ctxKey)
	if t, ok := v.(time.Time); ok {
		h.log.With(ctx).Info("redis-info",
			zap.String("cmd", cmd.String()),
			zap.String("runtime", time.Since(t).String()),
		)
	}
	return nil
}

//
// BeforeProcessPipeline 执行pipeline 之前
// @receiver h
// @param ctx
// @param cmds
// @return context.Context
// @return error
//
func (h Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	t := time.Now()
	ctx = context.WithValue(ctx, ctxKey, t)
	return ctx, nil
}

//
// AfterProcessPipeline 执行pipeline 之后
// @receiver h
// @param ctx
// @param cmds
// @return error
//
func (h Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	v := ctx.Value(ctxKey)
	if t, ok := v.(time.Time); ok {
		buf := new(bytes.Buffer)
		for _, v := range cmds {
			buf.WriteString(v.String())
			buf.WriteString(";")
		}
		h.log.With(ctx).Info("redis-info",
			zap.String("cmd", buf.String()),
			zap.String("runtime", time.Since(t).String()),
		)
	}
	return nil
}
