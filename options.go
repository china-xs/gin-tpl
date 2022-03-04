// Package gin_tpl
// @author: xs
// @date: 2022/3/3
// @Description: gin_tpl,描述
package gin_tpl

import "github.com/china-xs/gin-tpl/middleware"

type (
	Option  func(o *options)
	options struct {
		//ctx context.Context
		filters []FilterFunc
		ms      []middleware.Middleware
		port    int32 // 地址端口
	}
)

func Middleware(m ...middleware.Middleware) Option {
	return func(o *options) {
		o.ms = m
	}
}

func Filter(filters ...FilterFunc) Option {
	return func(o *options) {
		o.filters = filters
	}
}

func Port(port int32) Option {
	return func(o *options) {
		o.port = port
	}
}
