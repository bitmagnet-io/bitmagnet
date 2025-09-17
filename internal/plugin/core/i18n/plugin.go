package i18n

import (
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"go.uber.org/fx"
)

type deps struct{}

var (
	Ref = core.Ref.MustSub("i18n")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides translations for backend services"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithCliCommand[deps](NewCommand),
		builder.WithFxOption[deps](
			fx.Provide(
				i18n.NewBundle,
				// fx.Annotate(
				// 	func(providers []i18n.Provider) i18n.Provider {
				// 		return i18n.Providers(providers...)
				// 	},
				// 	fx.ParamTags(`group:"i18n_providers"`),
				// ),
			),
		),
	)
)
