package root

import (
	"github.com/bitmagnet-io/bitmagnet/internal/banner"
	"github.com/bitmagnet-io/bitmagnet/internal/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

func NewFactpry(builder *registry.Builder) cmd.CommandFactory {
	return func() cmd.Command {
		return &RootCommand{
			builder: builder,
		}
	}
}

type RootCommand struct {
	cmd.Cmd         `cmd:"name=bitmagnet"`
	builder         *registry.Builder
	Plugins         cmd.CSV[ref.Ref]
	DisabledPlugins cmd.CSV[ref.Ref]
	NoBanner        bool
	registry        *registry.Registry
}

func (cmd *RootCommand) Setup(env env.Env) error {
	if !cmd.NoBanner {
		banner.Write(env)
	}

	registry, err := cmd.builder.Resolve(
		env,
		registry.WithEnabledPlugins(cmd.Plugins...),
		registry.WithDisabledPlugins(cmd.DisabledPlugins...),
	)
	if err != nil {
		return err
	}

	cmd.registry = registry

	return nil
}

func (cmd *RootCommand) Subcommands() []cmd.Command {
	return cmd.registry.Commands()
}
