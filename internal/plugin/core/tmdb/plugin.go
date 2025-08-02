package tmdb

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb/tmdbhealthcheck"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	Config tmdb.Config

	deps struct {
		fx.In
		Client tmdb.Client
	}
)

var (
	Ref = core.Ref.MustSub("tmdb")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[Config, deps](),
		builder.WithDependencies[Config, deps](
			config.Ref,
			logging.Ref,
		),
		builder.WithDefaultConfig[Config, deps](Config(tmdb.NewDefaultConfig())),
		builder.WithFxOption[Config, deps](
			fx.Decorate(func(cfg Config, logger *zap.Logger, _ tmdb.Client) tmdb.Client {
				return tmdb.New(tmdb.Config(cfg), logger.Named(Ref.String()))
			}),
			fx.Decorate(func(tmdb.Enabled) tmdb.Enabled {
				return true
			}),
			fx.Decorate(func(_ tmdb.Config, cfg Config) tmdb.Config {
				return tmdb.Config(cfg)
			}),
		),
		builder.WithHealthCheckerOption[Config, deps](
			func(cfg Config, deps deps) health.CheckerOption {
				return tmdbhealthcheck.New(Ref.String(), deps.Client)
			},
		),
	)
)
