package recovery

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref = http_server.Ref.MustSub("recovery")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[deps](),
		builder.WithGinOption(Ref, func(deps) gin.OptionFunc {
			return func(engine *gin.Engine) {
				engine.Use(gin.Recovery())
			}
		}),
	)
)
