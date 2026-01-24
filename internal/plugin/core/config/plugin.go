package config

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config"
	"github.com/bitmagnet-io/bitmagnet/internal/config/lookup"
	"github.com/bitmagnet-io/bitmagnet/internal/config/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/spf13/afero"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref           = ref.Root.MustSub("config")
	RefActivation = Ref.MustSub("plugin_activation")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides configuration functionality"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
			fx.Provide(
				lookup.NewFromEnv,
				resolver.New,
				func(provider fs.Provider) manager.FS {
					return afero.Afero{
						Fs: afero.NewBasePathFs(provider.FSData(), config.SubpathPersisted),
					}
				},
				fx.Private,
			),
			fx.Provide(
				manager.New,
			),
		),
		builder.WithI18nMessage[deps](
			RefActivation,
			"description for plugin activation param",
			i18n.WithOther("Activation"),
		),
	)
)
