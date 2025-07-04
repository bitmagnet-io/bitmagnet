package tmdbfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb/tmdbhealthcheck"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"tmdb",
		configfx.NewConfigModule[tmdb.Config]("tmdb", tmdb.NewDefaultConfig()),
		fx.Provide(
			tmdb.New,
			fx.Annotate(
				func(client tmdb.Client, config tmdb.Config) health.CheckerOption {
					return tmdbhealthcheck.New(config.Enabled, client)
				},
				fx.ResultTags(`group:"health_check_options"`),
			),
		),
	)
}
