// Package blog
// @author: xs
// @date: 2022/3/7
// @Description: blog cmd 启动入口
// config viper

package main

import (
	tpl "github.com/china-xs/gin-tpl"
)

func main() {
	//tpl.
	app := tpl.NewServer()
	route, fc, err := initApp()
	// 初始化 路由
	route.InitRoute(app)
	if err != nil {
		panic(err)
	}
	defer func() {
		fc()
	}()
	if err := app.Run(); err != nil {
		panic(err)
	}

}
