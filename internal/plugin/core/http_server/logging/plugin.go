package logging

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/ginzap"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	config struct{}

	deps struct {
		fx.In
		Logger *zap.Logger
	}
)

var (
	Ref = http_server.Ref.MustSub("logging")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithDependencies[config, deps](
			logging.Ref,
		),
		builder.WithGinOption(
			Ref,
			func(_ config, deps deps) gin.OptionFunc {
				return func(engine *gin.Engine) {
					engine.Use(ginzap.Ginzap(deps.Logger, time.RFC3339, true))
				}
			},
		),
	)
)
