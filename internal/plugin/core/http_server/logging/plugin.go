package logging

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/httpserver/ginzap"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/logging"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
	Logger *zap.Logger
}

var (
	Ref = http_server.Ref.MustSub("logging")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides logging for the HTTP server"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			logging.Ref,
		),
		builder.WithGinOption(
			Ref,
			0,
			func(deps deps) gin.OptionFunc {
				return func(engine *gin.Engine) {
					engine.Use(ginzap.Ginzap(deps.Logger, time.RFC3339, true))
				}
			},
		),
	)
)
