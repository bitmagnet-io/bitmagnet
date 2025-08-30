package config

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/lookup"
	"github.com/bitmagnet-io/bitmagnet/internal/config/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/fs"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/spf13/afero"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref           = core.Ref.MustSub("config")
	RefActivation = Ref.MustSub("plugin_activation")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides configuration functionality"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
			fx.Provide(
				lookup.NewFromEnv,
				resolver.New,
				func(provider fs.FSProvider) manager.FS {
					return afero.Afero{
						Fs: afero.NewBasePathFs(provider.FSData(), config.SubpathPersisted),
					}
				},
				fx.Private,
			),
			fx.Provide(
				manager.New,
			),
			fx.Supply(
				fx.Annotate(
					&i18n.Message{
						ID:          RefActivation.String(),
						Description: "description for plugin activation param",
						Other:       "Activation",
					},
					fx.ResultTags(`group:"i18n_messages"`),
				),
			),
		),
		// builder.WithCliCommand[deps](
		// 	NewConfigCommand(),
		// ),
		// builder.WithFxOption[config, deps](configfx.New()),
	)
)
