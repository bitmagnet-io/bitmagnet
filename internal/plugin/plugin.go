package plugin

import (
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"go.uber.org/fx"
)

type ParamProvider interface {
	Params() []config_registry.Param
}

type Command func(app.Builder) cmd.Command

type CommandProvider interface {
	Commands() []Command
}

type FXOptionProvider interface {
	FXOption() fx.Option
}

type Plugin interface {
	Ref() ref.Ref
	ActivationRef() ref.Nullable
	Dependencies() []ref.Ref
	ParamProvider
	CommandProvider
	FXOptionProvider
}

func NewPlugin(
	ref ref.Ref,
	activation ref.Nullable,
	dependencies []ref.Ref,
	configParams []config_registry.Param,
	cliCommands []Command,
	fxOption fx.Option,
) Plugin {
	return &plugin{
		ref:          ref,
		activation:   activation,
		dependencies: dependencies,
		params:       configParams,
		cliCcommands: cliCommands,
		fxOption:     fxOption,
	}
}

type plugin struct {
	ref          ref.Ref
	activation   ref.Nullable
	dependencies []ref.Ref
	params       []config_registry.Param
	cliCcommands []Command
	fxOption     fx.Option
}

func (p *plugin) Ref() ref.Ref {
	return p.ref
}

func (p *plugin) ActivationRef() ref.Nullable {
	return p.activation
}

func (p *plugin) Dependencies() []ref.Ref {
	return p.dependencies
}

func (p *plugin) Params() []config_registry.Param {
	return p.params
}

func (p *plugin) FXOption() fx.Option {
	return fx.Module(
		p.Ref().String(),
		p.fxOption,
	)
}

func (p *plugin) Commands() []Command {
	if len(p.cliCcommands) == 0 {
		return nil
	}

	return p.cliCcommands
}
