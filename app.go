// Package gin_tpl
// @author: xs
// @date: 2022/3/3
// @Description: gin_tpl 目的，仅用于protobuf+gin合成配套服务
package gin_tpl

import (
	"context"
	"github.com/china-xs/gin-tpl/middleware"
	"github.com/gin-gonic/gin"
)

type (
	//AppInfo 服务发现注册信息，讲道理，不应该注册到服务器，也不应该提供grpc 如果需要直接上kratos
	AppInfo interface {
		ID() string
		Name() string
		Version() string
		Metadata() map[string]string
	}

	App struct {
		opts   options
		ctx    context.Context
		engine *gin.Engine

		filters []FilterFunc            //
		ms      []middleware.Middleware //
		cancel  func()
	}
)

// New 创建gin服务
func New(opts ...Option) *App {
	r := gin.Default()

	return &App{
		engine: r,
	}
}

//Run 启动
func (a *App) Run() error {

	a.engine.Run()
	return nil
}

// Stop 停止
func (a *App) Stop() error {
	return nil
}
