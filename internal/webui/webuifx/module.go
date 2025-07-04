package webuifx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/webui"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func New() fx.Option {
	return fx.Module(
		webui.Namespace,
		fx.Provide(
			fx.Annotate(
				func(logger *zap.SugaredLogger) httpserver.Option {
					return webui.New(logger.Named(webui.Namespace))
				},
				fx.ResultTags(`group:"http_server_options"`)),
		),
	)
}
