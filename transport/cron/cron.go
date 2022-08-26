// Package cron
// @author: xs
// @date: 2022/8/23
// @Description: cron 目的统计管理服务进程,只要实现 start&&stop 方法即可
package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
)

type ServerOption func(*Server)

func AddJob(sepc string, cmd cron.Job) ServerOption {
	return func(o *Server) {
		if _, err := o.c.AddJob(sepc, cmd); err != nil {
			panic(err)
		}
	}
}
func AddFunc(spec string, cmd func()) ServerOption {
	return func(o *Server) {
		if _, err := o.c.AddFunc(spec, cmd); err != nil {
			panic(err)
		}
	}
}

type Server struct {
	c *cron.Cron
}

func New(c *cron.Cron, opts ...ServerOption) *Server {
	s := &Server{
		c: c,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (c Server) Start(ctx context.Context) error {
	c.c.Start()
	// 接收到 统一取消信号
	<-ctx.Done()
	return nil
}

func (c Server) Stop(ctx context.Context) error {
	log.Print("[CRON] server stopping")
	ctx = c.c.Stop()
	// 定时任务结束
	<-ctx.Done()
	return nil
}
