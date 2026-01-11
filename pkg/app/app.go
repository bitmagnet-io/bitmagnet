package app

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd/root"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/provider"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

type App struct {
	bundles []plugin.Provider
}

func New(bundles ...plugin.Provider) *App {
	return &App{
		bundles: append(
			[]plugin.Provider{provider.Core()},
			bundles...,
		),
	}
}

func (a *App) Run(env env.Env) (int, error) {
	return root.NewFactpry(a.bundles...).Run(env)
}
