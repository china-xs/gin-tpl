//go:build wireinject
// +build wireinject

package main

import (
	"github.com/china-xs/gin-tpl/examples/blog/internal/server"
	"github.com/china-xs/gin-tpl/examples/blog/internal/service"
	"github.com/china-xs/gin-tpl/pkg/log"
	//tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/pkg/config"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	config.ProviderSet,
	server.InitRouteSet,
	service.ProviderSet,
)

// cf config path
func initApp(path string) (*server.Route, func(), error) {
	panic(wire.Build(providerSet))
}
