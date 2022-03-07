// Package main
// @author: xs
// @date: 2022/3/7
// @Description: main,描述
package main

import (
	tpl "github.com/china-xs/gin-tpl"
	route "github.com/china-xs/gin-tpl/example/blog/internal/server"
	"github.com/google/wire"
)

func initApp() (*tpl.Server, func(), error) {
	panic(
		wire.Build(
			route.Set, //路由注册
		))
}
