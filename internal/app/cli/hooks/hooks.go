package hooks

import (
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Hooks     []fx.Hook `group:"app_hooks"`
}

type Result struct {
	fx.Out
	AttachedHooks
}

type AttachedHooks struct{}

func New(p Params) (Result, error) {
	for _, hook := range p.Hooks {
		p.Lifecycle.Append(hook)
	}
	return Result{}, nil
}
