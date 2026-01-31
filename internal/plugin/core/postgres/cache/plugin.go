package cache

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/database/cache"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	caches "github.com/mgdigital/gorm-cache/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	libgorm "gorm.io/gorm"
)

type deps struct {
	fx.In
	TTL     *atomic.Value[cache.TTL]
	MaxKeys *atomic.Value[cache.MaxItems]
	Logger  *zap.Logger
}

var (
	Ref = postgres.Ref.MustSub("cache")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides Postgres query result caching"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithDependencies[deps](
			config.Ref,
		),
		builder.WithConfig[deps](Ref.MustSub("ttl"), cache.ParamTTL),
		builder.WithConfig[deps](Ref.MustSub("max_items"), cache.ParamMaxItems),
		builder.WithFxOption[deps](fx.Provide(cache.New)),
		builder.WithGormPlugin(
			func(deps deps) libgorm.Plugin {
				return &caches.Caches{Conf: &caches.Config{
					Cacher: cache.New(
						deps.TTL,
						deps.MaxKeys,
						deps.Logger.Named(Ref.String()),
					),
				}}
			},
		),
	)
)
