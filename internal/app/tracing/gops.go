package tracing

import (
	"context"
	"github.com/google/gops/agent"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	AppHook fx.Hook `group:"app_hooks"`
}

func New() Result {
	return Result{
		AppHook: fx.Hook{
			OnStart: func(ctx context.Context) error {
				return agent.Listen(agent.Options{})
			},
			OnStop: func(ctx context.Context) error {
				agent.Close()
				return nil
			},
		},
	}
}
