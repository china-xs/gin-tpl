// Package gin_tpl
// @author: xs
// @date: 2022/3/3
// @Description: gin_tpl 目的，仅用于protobuf+gin合成配套服务
package gin_tpl

import (
	"context"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/swagger-api/openapiv2"
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

//type siddleware(middleware.Handler) middleware.Handler

type Server struct {
	port    int32
	Engine  *gin.Engine
	timeout time.Duration
	ms      []middleware.Middleware
	Enc     EncodeResponseFunc
	sigs    []os.Signal
	srv     *http.Server
	apiSrv  *http.Server
	ctx     context.Context
	openApi bool // 是否开启接口文档
}

//OpenApi 是否开启文档
func OpenApi(b bool) ServerOption {
	return func(s *Server) {
		s.openApi = b
	}
}

// Timeout with server timeout. 服务超时时间，暂时没有控制请求期间超时，仅仅作用在热更新，延迟推出处理
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

//
// NewServer
// @Description: gin 启动器
// @param opts
// @return *Server
//
func NewServer(opts ...ServerOption) *Server {
	r := gin.Default()
	srv := &Server{
		Engine:  r,
		port:    8080,
		Enc:     DefaultResponseEncoder,
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		timeout: 10 * time.Second,
		ctx:     context.TODO(),
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
	//s.Engine.Handlers
	srv := &http.Server{
		//Addr:    fmt.Sprintf(":%v", s.port),
		Handler: s.Engine,
		Addr:    ":8080",
	}
	s.srv = srv
	eg.Go(func() error { return srv.ListenAndServe() })
	if s.openApi {
		eg.Go(func() error {
			openAPIhandler := openapiv2.NewHandler()
			swa := &http.Server{
				Addr:    ":8081",
				Handler: openAPIhandler,
			}
			s.apiSrv = swa
			log.Printf("open http://127.0.0.1:8081/q/swagger-ui#/")
			return swa.ListenAndServe()
		})
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, s.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Shutdown Server err ...")
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
	log.Printf("timeout:%v\n", s.timeout)
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	if s.apiSrv != nil {
		s.apiSrv.Shutdown(context.TODO())
	}
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

func (s *Server) Middleware(h middleware.Handler) middleware.Handler {
	return middleware.Chain(s.ms...)(h)
}
