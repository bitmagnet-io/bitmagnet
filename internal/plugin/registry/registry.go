package registry

import (
	"github.com/bitmagnet-io/bitmagnet/internal/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"go.uber.org/fx"
)

type Registry struct {
	pluginInfos []PluginInfo
	commands    []plugin.Command
	fxOption    fx.Option
}

func (r *Registry) PluginInfos() []PluginInfo {
	return append([]PluginInfo(nil), r.pluginInfos...)
}

func (r *Registry) Commands() []cmd.Command {
	return slice.Map(r.commands, func(cmd plugin.Command) cmd.Command {
		return cmd(r)
	})
}

func (r *Registry) Build(options ...fx.Option) *fx.App {
	return fx.New(append([]fx.Option{r.fxOption}, options...)...)
}
