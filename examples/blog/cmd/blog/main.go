// Package blog
// @author: xs
// @date: 2022/3/7
// @Description: blog cmd 启动入口
// config viper

package main

import (
	tpl "github.com/china-xs/gin-tpl"
	"github.com/china-xs/gin-tpl/middleware/validate"
	"time"
)

func main() {
	//tpl.
	var ops []tpl.ServerOption
	ms := tpl.Middleware(
		validate.Validator(),
	)

	ops = append(ops,
		ms,                 // 中间件
		tpl.OpenApi(false), //在线文档
		tpl.Timeout(5*time.Second),
	)
	app := tpl.NewServer(ops...)

	route, fc, err := initApp()
	// 初始化 路由
	route.InitRoute(app)
	if err != nil {
		panic(err)
	}
	defer fc()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
