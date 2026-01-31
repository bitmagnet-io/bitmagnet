package tmdb

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	deps struct {
		fx.In
		Client tmdb.Client
	}
)

var (
	Ref = ref.Root.MustSub("tmdb")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDependencies[deps](
			config.Ref,
			logging.Ref,
		),
		builder.WithActivation[deps](plugin.ActivationDisabled),
		builder.WithFxOption[deps](
			fx.Supply(tmdb.Enabled(false)),
			fx.Provide(func(logger *zap.Logger) tmdb.Client {
				return tmdb.New(tmdb.NewDefaultConfig(), logger.Named(Ref.String()))
			}),
			fx.Decorate(func(tmdb.Enabled) tmdb.Enabled {
				return false
			}),
		),
		// builder.WithHealthCheckerOption[Config, deps](
		// 	func(cfg Config, deps deps) health.CheckerOption {
		// 		return tmdbhealthcheck.New(Ref.String(), deps.Client)
		// 	},
		// ),
	)
)
