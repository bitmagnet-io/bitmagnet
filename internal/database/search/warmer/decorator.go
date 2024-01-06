package warmer

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DecoratorParams struct {
	fx.In
	Config   Config
	Search   lazy.Lazy[search.Search]
	Registry worker.Registry
	Logger   *zap.SugaredLogger
}

type DecoratorResult struct {
	fx.Out
	Registry worker.Registry
}

func NewDecorator(params DecoratorParams) (DecoratorResult, error) {
	var w warmer
	err := params.Registry.Decorate("http_server", func(hook fx.Hook) fx.Hook {
		return fx.Hook{
			OnStart: func(ctx context.Context) error {
				s, err := params.Search.Get()
				if err != nil {
					return err
				}
				w = warmer{
					stopped:  make(chan struct{}),
					interval: params.Config.Interval,
					search:   s,
					logger:   params.Logger.Named("search_warmer"),
				}
				go w.start()
				return hook.OnStart(ctx)
			},
			OnStop: func(ctx context.Context) error {
				if w.stopped != nil {
					close(w.stopped)
				}
				return hook.OnStop(ctx)
			},
		}
	})
	return DecoratorResult{
		Registry: params.Registry,
	}, err
}
