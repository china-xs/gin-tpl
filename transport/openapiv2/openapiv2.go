// Package openapiv2
// @author: xs
// @date: 2022/8/22
// @Description: openapiv2 swagger 文档方法为原始 http.Handler 与gin gin.Handler 不一致，
// 随独立开启新端口服务,服务遵循插拔式，可按照自己需求重造轮子
package openapiv2

import (
	"context"
	"github.com/china-xs/gin-tpl/errors"
	"github.com/go-kratos/swagger-api/openapiv2"
	"log"
	"net/http"
)

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

type Server struct {
	*http.Server
	addr string
}

func Addr(addr string) ServerOption {
	return func(o *Server) {
		o.addr = addr
	}
}

func New(opts ...ServerOption) *Server {
	o := &Server{
		addr: "0.0.0.0:8081",
	}
	for _, v := range opts {
		v(o)
	}
	return o
}

func (s *Server) Start(ctx context.Context) error {
	h := openapiv2.NewHandler()
	log.Printf("listen:%v\n", s.addr)
	s.Server = &http.Server{
		Addr:    s.addr,
		Handler: h,
	}
	err := s.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	log.Print("[SWAGGER] server stopping")
	err := s.Shutdown(ctx)
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
