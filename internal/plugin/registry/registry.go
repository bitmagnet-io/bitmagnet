package registry

import (
	"github.com/bitmagnet-io/bitmagnet/internal/cmd"
	config_resolver "github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type (
	Registry struct {
		config
		instances ref.Map[plugin.Instance]
		commands  []plugin.Command
		fxOption  fx.Option
	}

	config = config_resolver.Resolved
)

func (r *Registry) Instances() ref.Map[plugin.Instance] {
	return r.instances.Clone()
}

func (r *Registry) Commands() []cmd.Command {
	return slice.Map(r.commands, func(cmd plugin.Command) cmd.Command {
		return cmd(r)
	})
}

func (r *Registry) Build(options ...fx.Option) *fx.App {
	return fx.New(append([]fx.Option{r.fxOption}, options...)...)
}
