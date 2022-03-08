//go:build wireinject
// +build wireinject

package main

import (
	//tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/examples/blog/internal/server"
	"github.com/china-xs/gin-tpl/examples/blog/internal/service"
	"github.com/google/wire"
)

// cf config path
func initApp() (*server.Route, func(), error) {
	panic(wire.Build(
		server.InitRouteSet,
		service.Set))
}
