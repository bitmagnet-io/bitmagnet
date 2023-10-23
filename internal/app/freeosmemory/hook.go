package freeosmemory

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

type Params struct {
	fx.In
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	AppHook fx.Hook `group:"app_hooks"`
}

func New(p Params) Result {
	logger := p.Logger.Named("free_os_memory")
	var done chan struct{}
	return Result{
		AppHook: fx.Hook{
			OnStart: func(context.Context) error {
				done = make(chan struct{})
				go (func() {
					for {
						select {
						case <-done:
							return
						case <-time.After(5 * time.Minute):
							logger.Debug("freeing os memory")
							debug.FreeOSMemory()
						}
					}
				})()
				return nil
			},
			OnStop: func(context.Context) error {
				close(done)
				return nil
			},
		},
	}
}
