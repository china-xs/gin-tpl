// Package blog
// @author: xs
// @date: 2022/3/7
// @Description: blog cmd 启动入口
// config viper

package main

import (
	"flag"
	tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/examples/blog/internal/server"
	"github.com/china-xs/gin-tpl/middleware/logger"
	"github.com/china-xs/gin-tpl/middleware/validate"
	"github.com/kataras/i18n"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// go build -ldflags "-X main.Version=x.y.z"
//var (
//	// flagconf is the config flag.
//	flagconf string
//)
//
//func init() {
//	flag.StringVar(&flagconf, "conf", "../../configs/app.yaml", "config path, eg: -conf config.yaml")
//}

var configFile = flag.String("f", "../../configs/app.yaml", "set config file which viper will loading.")

func main() {
	flag.Parse()
	app, fc, err := initApp(*configFile)
	if err != nil {
		panic(err)
	}
	defer fc()
	if err := app.Run(); err != nil {
		panic(err)
	}

}

func newApp(route server.Route, log *zap.Logger, v *viper.Viper) *tpl.Server {
	var ops []tpl.ServerOption
	I18n, err := i18n.New(i18n.Glob("../../configs/locales/*/*"), "en-US", "zh-CN")
	if err != nil {
		panic(err)
	}
	ms := tpl.Middleware(
		validate.Validator2I18n(I18n),
		logger.Logger(log),
	)
	opts, err := tpl.NewSerOpts(v)
	if err != nil {
		panic(err)
	}
	ops = append(ops, ms, opts)
	app := tpl.NewServer(ops...)
	route.InitRoute(app)
	return app
}
