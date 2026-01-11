package plugin

import (
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"go.uber.org/fx"
)

type ConfigProvider interface {
	ConfigParams() []config_registry.Param
}

type ErrorProvider interface {
	Errors() ref.Map[error]
}

type I18nProvider interface {
	I18nMessages() []*i18n.Message
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
	ConfigProvider
	ErrorProvider
	I18nProvider
	CommandProvider
	FXOptionProvider
}

func NewPlugin(
	ref ref.Ref,
	activation ref.Nullable,
	dependencies []ref.Ref,
	configParams []config_registry.Param,
	errors ref.Map[error],
	i18nProvider i18n.Provider,
	cliCommands []Command,
	fxOption fx.Option,
) Plugin {
	return &plugin{
		ref:          ref,
		activation:   activation,
		dependencies: dependencies,
		configParams: configParams,
		errors:       errors,
		i18nProvider: i18nProvider,
		cliCcommands: cliCommands,
		fxOption:     fxOption,
	}
}

type plugin struct {
	ref          ref.Ref
	activation   ref.Nullable
	dependencies []ref.Ref
	configParams []config_registry.Param
	errors       ref.Map[error]
	i18nProvider i18n.Provider
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

func (p *plugin) ConfigParams() []config_registry.Param {
	return p.configParams
}

func (p *plugin) Errors() ref.Map[error] {
	return p.errors
}

func (p *plugin) I18nMessages() []*i18n.Message {
	return p.i18nProvider()
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
