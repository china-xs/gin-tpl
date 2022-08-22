// Package http
// @author: xs
// @date: 2022/7/22
// @Description: http
package http

import (
	"context"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
)

const OperationKey = "operation"

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

type Server struct {
	*http.Server
	addr   string // default 0.0.0.0:8080
	Engine *gin.Engine
	dec    DecodeRequestFunc       // 请求参数绑定结构
	enc    EncodeResponseFunc      // 定义返回结构
	ms     []middleware.Middleware // 全局中间价
}

// Addr with service addr option.
func Addr(addr string) ServerOption {
	return func(o *Server) {
		o.addr = addr
	}
}

// Middleware with service middleware option.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.ms = m
	}
}

// RequestDecoder with request decoder.
func RequestDecoder(dec DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.dec = dec
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(en EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.enc = en
	}
}

// NewServer creates an HTTP server by options.
func NewServer(opts ...ServerOption) *Server {
	r := gin.Default()
	srv := &Server{
		Engine: r,
		addr:   "0.0.0.0:8080",
		dec:    DefaultRequestDecoder,
		enc:    DefaultResponseEncoder,
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

//
// Start this HTTP server.
// @receiver s
// @param ctx
// @return error
//
func (s *Server) Start(ctx context.Context) error {
	// ReadTimeout 超时时间只存在ctx，不会有程序上退出
	log.Printf("listen:%v\n", s.addr)
	// 超时控制可以按照这个实现
	// https://toutiao.io/posts/ytrjt0/preview
	s.Server = &http.Server{
		Addr:    s.addr,
		Handler: s.Engine,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return ctx
		},
	}
	err := s.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	// 暂时不加tls
	return err
}

//
// Stop the HTTP server.
// @receiver s
// @param ctx
// @return error
//
func (s *Server) Stop(ctx context.Context) error {
	log.Print("[HTTP] server stopping")
	err := s.Shutdown(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

//
// Middleware gin 绑定实现方法
// @receiver s
//
func (s *Server) Middleware(h middleware.Handler) middleware.Handler {
	return middleware.Chain(s.ms...)(h)
}

// Bind  请求参数绑定
func (s *Server) Bind(c *gin.Context, obj interface{}) error {
	return s.dec(c, obj)
}

// Result 结果结返回
func (s *Server) Result(c *gin.Context, obj interface{}, err error) {
	s.enc(c, obj, err)
	return
}
