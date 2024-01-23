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
	Config Config
	Search lazy.Lazy[search.Search]
	Logger *zap.SugaredLogger
}

type DecoratorResult struct {
	fx.Out
	Decorator worker.Decorator `group:"worker_decorators"`
}

func New(params DecoratorParams) DecoratorResult {
	var w warmer
	return DecoratorResult{
		Decorator: worker.Decorator{
			Key: "http_server",
			Decorate: func(hook fx.Hook) fx.Hook {
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
			},
		},
	}
}
