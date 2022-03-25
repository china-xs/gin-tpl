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
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const OperationKey = "operation"

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

type Server struct {
	port    int32 // 端口
	Engine  *gin.Engine
	timeout time.Duration           // 请求超时时长
	ms      []middleware.Middleware // 中间价
	filters []gin.HandlerFunc       // gin 中间件
	Enc     EncodeResponseFunc
	sigs    []os.Signal
	srv     *http.Server
	apiSrv  *http.Server
	ctx     context.Context
	openApi bool // 是否开启接口文档
	name    string
}

type Opts struct {
	Name    string        `yaml:"name"`
	Port    int32         `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
	OpenApi bool          `yaml:"openApi"`
}

func NewSerOpts(v *viper.Viper) (serv ServerOption, err error) {
	o := new(Opts)
	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}
	return func(s *Server) {
		s.name = o.Name
		s.port = o.Port
		s.timeout = o.Timeout
		s.openApi = o.OpenApi
	}, nil
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

//
// Middleware with service middleware option.
// @Description: 执行方法过滤器
// @param m
// @return ServerOption
//
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.ms = m
	}
}

//
// Name
// @Description: 服务名称
// @param name
// @return ServerOption
//
func Name(name string) ServerOption {
	return func(o *Server) {
		o.name = name
	}
}

//
// ResponseEncoder with response encoder.
// @Description: 返回前端封装方法
// @param en
// @return ServerOption
//
func ResponseEncoder(en EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.Enc = en
	}
}

//
// Port
// @Description: 服务监听端口
// @param port
// @return ServerOption
//
func Port(port int32) ServerOption {
	return func(o *Server) {
		o.port = port
	}
}

//
// Filter
// @Description: gin 全局中间件，无法覆盖链路中间件
// @param filters
// @return ServerOption
//
func Filter(filters ...gin.HandlerFunc) ServerOption {
	return func(o *Server) {
		o.filters = filters
	}
}

//
// Signal with exit signals.
// @Description: 热重启信号
// @param sigs
// @return ServerOption
//
func Signal(sigs ...os.Signal) ServerOption {
	return func(o *Server) { o.sigs = sigs }
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
		name:    "gin-app",
	}
	for _, o := range opts {
		o(srv)
	}
	// 全局中间件
	if len(srv.filters) > 0 {
		r.Use(srv.filters...)
	}
	// 链路全局注册
	tp := traceSDK.NewTracerProvider(
		traceSDK.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(srv.name),
		)),
	)
	r.Use(otelgin.Middleware(srv.name))
	//srv.Engine = r
	//Tracer 全局注册
	otel.SetTracerProvider(tp)
	return srv
}

//Run 启动
func (s *Server) Run() error {
	ctx := context.TODO()
	eg, ctx := errgroup.WithContext(ctx)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", s.port),
		Handler: s.Engine,
	}
	s.srv = srv
	eg.Go(func() error {
		log.Printf("listen:%v", s.port)
		return srv.ListenAndServe()
	})
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

func GetOperation(c *gin.Context) string {
	operation, _ := c.Get(OperationKey)
	return operation.(string)
}
