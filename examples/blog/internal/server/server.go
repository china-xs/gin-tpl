// Package server
// @author: xs
// @date: 2022/3/7
// @Description: 路由注册
package server

import (
	tpl "github.com/china-xs/gin-tpl"
	apiAuth "github.com/china-xs/gin-tpl/examples/blog/api/auth"
	implLogin "github.com/china-xs/gin-tpl/examples/blog/internal/service/auth"
	"github.com/google/wire"
)

//var Provider = wire.NewSet(NewRoute, NewOptions)

//wire.NewSet
var InitRouteSet = wire.NewSet(wire.Struct(new(Route), "*"))

type Route struct {
	SrvLogin *implLogin.LoginService
}

func (r Route) InitRoute(app *tpl.Server) (*tpl.Server, error) {
	apiAuth.RegisterLoginGinServer(app, r.SrvLogin)
	return app, nil
}

func registerGraph(app *tpl.Server) {
	app.Route("GET", "/graphql")

}
