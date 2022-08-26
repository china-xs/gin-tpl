// Package gin_tpl
// @author: xs
// @date: 2022/7/21
// @Description: gin_tpl
package gin_tpl

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Server 迭代启动关闭统一由CoreServer 管控
type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type App struct {
	ctx    context.Context
	cancel func()
	opts   options
}

func New(opts ...Option) *App {
	ctx := context.TODO()

	o := options{
		ctx:     context.TODO(),
		sigs:    []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		timeout: 10 * time.Second,
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

func (a App) Run() error {
	eg, ctx := errgroup.WithContext(a.ctx)
	var wg sync.WaitGroup
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done()
			ctx1 := context.TODO()
			stopCtx, cancel := context.WithTimeout(ctx1, a.opts.timeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			// 不可以使用defer
			wg.Done()
			return srv.Start(ctx)
		})
	}
	// 等待服务启动完成
	wg.Wait()
	c := make(chan os.Signal, 1)
	// 监听关闭进程信号
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})
	// 等待所有eg所有进程结束
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a App) Stop() error {
	// 执行发送cancel 信号
	a.cancel()
	return nil
}
