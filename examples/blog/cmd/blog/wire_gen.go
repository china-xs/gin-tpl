// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/china-xs/gin-tpl/examples/blog/internal/server"
	"github.com/china-xs/gin-tpl/examples/blog/internal/service"
	"github.com/china-xs/gin-tpl/examples/blog/internal/service/auth"
	"github.com/china-xs/gin-tpl/pkg/config"
	"github.com/china-xs/gin-tpl/pkg/db"
	"github.com/china-xs/gin-tpl/pkg/log"
	"github.com/china-xs/gin-tpl/pkg/redis"
	"github.com/google/wire"
)

// Injectors from wire.go:

// cf config path
func initApp(path string) (*server.Route, func(), error) {
	viper, err := config.New(path)
	if err != nil {
		return nil, nil, err
	}
	options, err := log.NewOptions(viper)
	if err != nil {
		return nil, nil, err
	}
	logger, err := log.New(options)
	if err != nil {
		return nil, nil, err
	}
	dbOptions, err := db.New(viper)
	if err != nil {
		return nil, nil, err
	}
	gormDB, cleanup, err := db.NewDb(dbOptions, logger)
	if err != nil {
		return nil, nil, err
	}
	redisOptions, err := redis.NewOps(viper)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	client, err := redis.New(redisOptions)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	loginService := auth.NewLoginService(logger, gormDB, client)
	route := &server.Route{
		SrvLogin: loginService,
	}
	return route, func() {
		cleanup()
	}, nil
}

// wire.go:

var providerSet = wire.NewSet(log.ProviderSet, db.ProviderSet, redis.ProviderSet, config.ProviderSet, server.InitRouteSet, service.ProviderSet)
