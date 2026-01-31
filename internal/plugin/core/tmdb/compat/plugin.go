package compat

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/tmdb"
	internaltmdb "github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref = tmdb.Ref.MustSub("compat")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Annotate(
					func() (client internaltmdb.Client, config internaltmdb.Config, enabled internaltmdb.Enabled) {
						return
					},
				),
			),
		),
	)
)
