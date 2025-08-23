package cors

import (
	"fmt"
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
	AllowedOrigins AllowedOrigins
	Logger         *zap.Logger
}

var (
	Ref = http_server.Ref.MustSub("cors")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[deps](),
		builder.WithDependencies[deps](
			config.Ref,
			logging.Ref,
		),
		// builder.WithDefaultConfig[deps](NewDefaultConfig()),
		builder.WithConfigParam[deps](Ref.MustSub("allowed_origins"), ParamAllowedOrigins),
		builder.WithGinOption(
			Ref,
			func(deps deps) gin.OptionFunc {
				return func(engine *gin.Engine) {
					engine.Use(gincors.New(cors.Options{
						AllowedOrigins: deps.AllowedOrigins,
						Debug:          true,
						// we don't need every request logged so apply sampling
						Logger: corsLogger{deps.Logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
							return zapcore.NewSamplerWithOptions(core, time.Hour, 10, 0)
						})).Named(Ref.String())},
					}))
				}
			}),
	)
)

type corsLogger struct {
	logger *zap.Logger
}

func (c corsLogger) Printf(format string, v ...interface{}) {
	c.logger.Debug(fmt.Sprintf(format, v...))
}
