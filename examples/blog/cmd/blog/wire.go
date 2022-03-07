//go:build wireinject
// +build wireinject

package main

import (
	tpl "github.com/china-xs/gin-tpl"
	route "github.com/china-xs/gin-tpl/example/blog/internal/server"
	srv "github.com/china-xs/gin-tpl/example/blog/internal/service"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
// pkg init
)

// cf config path
func initApp(cf string) (*tpl.Server, func(), error) {
	panic(
		wire.Build(
			// db init
			// log init
			route.Set, //路由注册
			srv.Set,   // 接口实现
		))
}
