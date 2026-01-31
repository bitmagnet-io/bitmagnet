package plugin

import (
	"context"

	embed_i18n "github.com/bitmagnet-io/bitmagnet/i18n"
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/pkg/i18n"
	pkg_plugin "github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type (
	Command    = pkg_plugin.Command
	Plugin     = pkg_plugin.Plugin
	Activation = pkg_plugin.Activation
)

const (
	ActivationEnabled  = pkg_plugin.ActivationEnabled
	ActivationDisabled = pkg_plugin.ActivationDisabled
	ActivationAuto     = pkg_plugin.ActivationAuto
	ActivationAlways   = pkg_plugin.ActivationAlways

	KeyActivation = pkg_plugin.KeyActivation
)

func NewPlugin(
	ref ref.Ref,
	activation ref.Nullable,
	dependencies []ref.Ref,
	configParams []config_registry.Param,
	errors ref.Map[error],
	i18nProvider i18n.MessageProvider,
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
	i18nProvider i18n.MessageProvider
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

func (p *plugin) LocalizedContent(
	_ context.Context,
	acceptLanguage ...string,
) (pkg_plugin.LocalizedContent, error) {
	localizer := i18n.NewLocalizer(embed_i18n.Bundle, acceptLanguage...)

	var description string
	if localized, _ := localizer.LocalizeMessage(&i18n.Message{
		ID: p.ref.String(),
	}); localized != "" {
		description = localized
	}

	return pkg_plugin.LocalizedContent{
		Ref:         p.ref,
		Description: description,
		ConfigParams: slice.Map(
			p.configParams,
			func(p config_registry.Param) pkg_plugin.LocalizedConfigParam {
				var label string
				if localized, _ := localizer.LocalizeMessage(&i18n.Message{
					ID: p.Ref.String(),
				}); localized != "" {
					label = localized
				} else {
					label = p.Description()
				}

				return pkg_plugin.LocalizedConfigParam{
					Ref:         p.Ref,
					Description: label,
				}
			},
		),
	}, nil
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
