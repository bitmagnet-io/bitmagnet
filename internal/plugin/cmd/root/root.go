package root

import (
	"bytes"
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/banner"
	"github.com/bitmagnet-io/bitmagnet/internal/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	wasm_plugin "github.com/bitmagnet-io/bitmagnet/internal/wasm/plugin"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

func NewFactpry(providers ...plugin.Provider) cmd.CommandFactory {
	return func() cmd.Command {
		return &Command{
			providers: providers,
		}
	}
}

// revive:disable:line-length-limit
type Command struct {
	cmd.Cmd       `cmd:"name=bitmagnet"`
	providers     []plugin.Provider
	LoadPlugin    cmd.CSV[loadPlugin] `cmd:"doc=Load plugins from the specified paths,example='/path/to/plugin,alias:/path/to/plugin'"`
	EnablePlugin  cmd.CSV[ref.Ref]    `cmd:"doc=Enable plugins,example=plugin"`
	DisablePlugin cmd.CSV[ref.Ref]    `cmd:"doc=Disable plugins,example=plugin"`
	NoBanner      bool                `cmd:"doc=Do not show the banner"`
	registry      *registry.Registry
}

func (cmd *Command) Setup(env env.Env) error {
	if !cmd.NoBanner {
		banner.Write(env)
	}

	providers := cmd.providers

	if len(cmd.LoadPlugin) > 0 {
		provider, err := wasm_plugin.NewProvider(slice.Map(
			cmd.LoadPlugin,
			func(lp loadPlugin) wasm_plugin.ProviderOption {
				return wasm_plugin.LoadPlugin(lp.path, lp.alias)
			})...,
		)
		if err != nil {
			return err
		}

		providers = append(
			providers,
			provider,
		)
	}

	registry, err := registry.New(providers...).Resolve(
		env,
		registry.WithEnabledPlugins(cmd.EnablePlugin...),
		registry.WithDisabledPlugins(cmd.DisablePlugin...),
	)
	if err != nil {
		return err
	}

	cmd.registry = registry

	return nil
}

func (cmd *Command) Subcommands() []cmd.Command {
	return cmd.registry.Commands()
}

type loadPlugin struct {
	path  string
	alias string
}

func (l *loadPlugin) UnmarshalText(text []byte) error {
	parts := bytes.SplitN(text, []byte("="), 2)
	switch len(parts) {
	case 1:
		l.path = string(parts[0])
	case 2:
		l.alias = string(parts[0])
		l.path = string(parts[1])
	}

	if l.path == "" {
		return errors.New("plugin path cannot be empty")
	}

	return nil
}
