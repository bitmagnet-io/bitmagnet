package app

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd/root"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/provider"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

type App struct {
	providers []plugin.Provider
}

func New(providers ...plugin.Provider) *App {
	return &App{
		providers: append(
			[]plugin.Provider{provider.Core()},
			providers...,
		),
	}
}

func (a *App) Run(env env.Env) (int, error) {
	return root.NewFactpry(a.providers...).Run(env)
}
