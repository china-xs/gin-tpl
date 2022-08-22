// Package gin_tpl
// @author: xs
// @date: 2022/7/29
// @Description: gin_tpl
package gin_tpl

import (
	"context"
	"os"
)

type Option func(o *options)

type options struct {
	ctx     context.Context
	sigs    []os.Signal
	servers []CoreServer
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

func Servers(srv ...CoreServer) Option {
	return func(o *options) { o.servers = srv }
}
