package cache

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/cache"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/gorm"
	caches "github.com/mgdigital/gorm-cache/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	libgorm "gorm.io/gorm"
)

type (
	Config = cache.Config

	deps struct {
		fx.In
		Logger *zap.SugaredLogger
	}
)

var (
	Ref = gorm.Ref.MustSub("cache")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[Config, deps](),
		builder.WithDependencies[Config, deps](
			config.Ref,
		),
		builder.WithDefaultConfig[Config, deps](cache.NewDefaultConfig()),
		builder.WithFxOption[Config, deps](fx.Provide(cache.New)),
		builder.WithGormPlugin(
			func(cfg Config, deps deps) libgorm.Plugin {
				return &caches.Caches{Conf: &caches.Config{
					Cacher: cache.New(
						cfg.TTL,
						int(cfg.MaxKeys),
						deps.Logger.Named(Ref.String()),
					),
				}}
			},
		),
	)
)
