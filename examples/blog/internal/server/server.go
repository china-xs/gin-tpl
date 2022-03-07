// Package server
// @author: xs
// @date: 2022/3/7
// @Description: 路由注册
package server

import (
	tpl "github.com/china-xs/gin-tpl"
	apiAuth "github.com/china-xs/gin-tpl/example/blog/api/auth"
	implLogin "github.com/china-xs/gin-tpl/example/blog/internal/service/auth"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewRoute,
)

var InitRouteSet = wire.NewSet(wire.Struct(new(RouteSet), "*"))

type RouteSet struct {
	srvLogin *implLogin.LoginService
	//ReportRepo *data.MpQu
	//ActivityInstancesRepo *merchants.ActivityInstancesRepo
}

//
func NewRoute(app *tpl.Server, routes RouteSet) (func(), error) {
	//注册路由
	apiAuth.RegisterLoginGinServer(app, routes.srvLogin)

	return func() {

	}, nil
}
