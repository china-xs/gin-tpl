// Package tasks
// @author: xs
// @date: 2022/8/26
// @Description: tasks
package tasks

import (
	"context"
	"github.com/china-xs/gin-tpl/errors"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"net/http"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

// RedisClick with set server driver
func RedisClick(rdb *redis.Client) ServerOption {
	return func(o *Server) {
		o.rdb = rdb
	}
}

// Config with asynq.config server 配置
func Config(c asynq.Config) ServerOption {
	return func(o *Server) {
		o.config = c
	}
}

// Mux with asynq.Server 路由 消费路由定义
func Mux(mux *asynq.ServeMux) ServerOption {
	return func(o *Server) {
		o.mux = mux
	}
}

type Server struct {
	rdb    *redis.Client
	config asynq.Config
	mux    *asynq.ServeMux
	srv    *asynq.Server
}

// New create asynq.Server
func New(opts ...ServerOption) *Server {
	srv := &Server{}
	for _, v := range opts {
		v(srv)
	}
	srv.srv = asynq.NewServer(srv, srv.config)
	return srv
}

// NewClient create asynq.Client
func NewClient(rdb *redis.Client) *asynq.Client {
	s := Server{rdb: rdb}
	return asynq.NewClient(s)
}

// MakeRedisClient redis 实现asynq  RedisConnOpt
func (s Server) MakeRedisClient() interface{} {
	return s.rdb
}

// Start 服务统一启动
func (s Server) Start(ctx context.Context) error {
	if s.mux == nil {
		return errors.New(http.StatusInternalServerError, "task mux is null", "TASK_MUX_ERR")
	}
	if s.rdb == nil {
		return errors.New(http.StatusBadRequest, "task redis.Client is null", "TASK_REDIS_ERR")
	}
	return s.srv.Start(s.mux)
}

// Stop 服务统一关闭
func (s Server) Stop(ctx context.Context) error {
	s.srv.Shutdown()
	s.srv.Stop()
	return nil
}
