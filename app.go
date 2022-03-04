// Package gin_tpl
// @author: xs
// @date: 2022/3/3
// @Description: gin_tpl 目的，仅用于protobuf+gin合成配套服务
package gin_tpl

import (
	"fmt"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
	"time"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

type Server struct {
	port    int32
	Engine  *gin.Engine
	timeout time.Duration
	ms      []middleware.Middleware
	Enc     EncodeResponseFunc
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Middleware with service middleware option.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.ms = m
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(en EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.Enc = en
	}
}

func NewServer(opts ...ServerOption) *Server {
	r := gin.Default()
	srv := &Server{
		Engine: r,
		port:   8080,
	}

	for _, o := range opts {
		o(srv)
	}
	return srv
}

//Run 启动
func (s *Server) Run() error {
	s.Engine.Run(fmt.Sprintf(":%v", s.port))
	return nil
}

// Stop 停止
func (s *Server) Stop() error {
	return nil
}

// Route 新增自定义路由  文件上传
func (s *Server) Route(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	s.Engine.Handle(httpMethod, relativePath, handlers...)
}
