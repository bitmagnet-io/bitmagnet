package redisfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/redis/healthcheck"
	"github.com/bitmagnet-io/bitmagnet/internal/redis/redisconfig"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"redis",
		configfx.NewConfigModule[redisconfig.Config]("redis", redisconfig.NewDefaultConfig()),
		fx.Provide(
			func(cfg redisconfig.Config) lazy.Lazy[*redis.Client] {
				return lazy.New(func() (*redis.Client, error) {
					return redis.NewClient(cfg.RedisClientOptions()), nil
				})
			},
			healthcheck.New,
		),
	)
}
