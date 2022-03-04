// Package gin_tpl
// @author: xs
// @date: 2022/3/3
// @Description: gin_tpl 目的，仅用于protobuf+gin合成配套服务
package gin_tpl

import (
	"context"
	"fmt"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	sigs    []os.Signal
	srv     *http.Server
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
		Enc:    DefaultResponseEncoder,
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}

	for _, o := range opts {
		o(srv)
	}
	return srv
}

//Run 启动
func (s *Server) Run() error {
	ctx := context.TODO()
	eg, ctx := errgroup.WithContext(ctx)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: s.Engine,
	}
	s.srv = srv
	eg.Go(func() error { return srv.ListenAndServe() })
	c := make(chan os.Signal, 1)
	signal.Notify(c, s.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				log.Printf("Shutdown Server ...")
				err := s.Stop()
				if err != nil {
					//s.opts.logger.Errorf("failed to stop app: %v", err)
					return err
				}
			}
		}
	})
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Stop 停止
func (s *Server) Stop() error {
	// 延迟关闭服务器
	ctx, cancel := context.WithTimeout(context.TODO(), s.timeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
		return err
	}
	log.Printf("Server exiting")
	return nil
}

// Route 新增自定义路由  文件上传
func (s *Server) Route(httpMethod, relativePath string, handlers ...gin.HandlerFunc) {
	s.Engine.Handle(httpMethod, relativePath, handlers...)
}
