package cache

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/database/cache"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/gorm"
	caches "github.com/mgdigital/gorm-cache/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	libgorm "gorm.io/gorm"
)

type deps struct {
	fx.In
	TTL     *atomic.Value[cache.TTL]
	MaxKeys *atomic.Value[cache.MaxKeys]
	Logger  *zap.Logger
}

var (
	Ref = gorm.Ref.MustSub("cache")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[deps](),
		builder.WithDependencies[deps](
			config.Ref,
		),
		builder.WithConfigParam[deps](Ref.MustSub("ttl"), cache.ParamTTL),
		builder.WithConfigParam[deps](Ref.MustSub("max_keys"), cache.ParamMaxKeys),
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
