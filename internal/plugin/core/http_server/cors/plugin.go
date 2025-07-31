package cors

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	gincors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type deps struct {
	fx.In
	Logger *zap.SugaredLogger
}

var (
	Ref = http_server.Ref.MustSub("cors")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[Config, deps](),
		builder.WithDependencies[Config, deps](
			config.Ref,
			logging.Ref,
		),
		builder.WithDefaultConfig[Config, deps](NewDefaultConfig()),
		builder.WithGinOption(
			Ref,
			func(cfg Config, deps deps) gin.OptionFunc {
				return func(engine *gin.Engine) {
					engine.Use(gincors.New(cors.Options{
						AllowedOrigins:       cfg.AllowedOrigins,
						AllowedMethods:       cfg.AllowedMethods,
						AllowedHeaders:       cfg.AllowedHeaders,
						ExposedHeaders:       cfg.ExposedHeaders,
						MaxAge:               cfg.MaxAge,
						AllowCredentials:     cfg.AllowCredentials,
						AllowPrivateNetwork:  cfg.AllowPrivateNetwork,
						OptionsPassthrough:   cfg.OptionsPassthrough,
						OptionsSuccessStatus: cfg.OptionsSuccessStatus,
						Debug:                cfg.Debug,
						// we don't need every request logged so apply sampling
						Logger: corsLogger{deps.Logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
							return zapcore.NewSamplerWithOptions(core, time.Hour, 10, 0)
						})).Named("cors")},
					}))
				}
			}),
	)
)

type corsLogger struct {
	logger *zap.SugaredLogger
}

func (c corsLogger) Printf(format string, v ...interface{}) {
	c.logger.Debugf(format, v...)
}
