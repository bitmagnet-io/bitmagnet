package plugin

import (
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
)

type Provider interface {
	LoadPlugins(env.Env) ([]Plugin, error)
}

type Plugins []Plugin

func (p Plugins) LoadPlugins(env.Env) ([]Plugin, error) {
	return p, nil
}
