package app

import (
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/bundle"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd/root"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/registry"
)

func Run(env env.Env, bundles ...bundle.Bundle) (int, error) {
	app, err := New(bundles...)
	if err != nil {
		return 1, err
	}

	return app.Run(env)
}

type App struct {
	registry *registry.Builder
}

func New(bundles ...bundle.Bundle) (*App, error) {
	coreBundle, err := bundle.Core()
	if err != nil {
		return nil, err
	}

	registry, err := registry.New(
		append(
			[]bundle.Bundle{coreBundle},
			bundles...,
		)...,
	)

	if err != nil {
		return nil, err
	}

	return &App{
		registry: registry,
	}, nil
}

func (k *App) Run(env env.Env) (int, error) {
	return root.NewFactpry(k.registry).Run(env)
}
