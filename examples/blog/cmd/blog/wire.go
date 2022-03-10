//go:build wireinject
// +build wireinject

package main

import (
	tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/examples/blog/internal/server"
	"github.com/china-xs/gin-tpl/examples/blog/internal/service"
	"github.com/china-xs/gin-tpl/pkg/db"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/china-xs/gin-tpl/pkg/redis"

	//tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/pkg/config"
	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	db.ProviderSet,
	redis.ProviderSet,
	config.ProviderSet,
	server.InitRouteSet,
	service.ProviderSet,
)

// cf config path
func initApp(path string) (*tpl.Server, func(), error) {
	panic(wire.Build(
		providerSet,
		newApp))
}
