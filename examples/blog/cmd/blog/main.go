// Package blog
// @author: xs
// @date: 2022/3/7
// @Description: blog cmd 启动入口
// config viper

package main

import (
	"flag"
	tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/middleware/validate"
	"time"
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
	//tpl.
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		validate.Validator(),
	)
	ops = append(ops,
		ms,                // 中间件
		tpl.OpenApi(true), //在线文档
		tpl.Timeout(5*time.Second),
		tpl.Name("gin-blog"),
		//tpl.Port(9090),
	)
	app := tpl.NewServer(ops...)

	route, fc, err := initApp(*configFile)
	if err != nil {
		panic(err)
	}
	// 初始化 路由
	route.InitRoute(app)
	defer fc()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
