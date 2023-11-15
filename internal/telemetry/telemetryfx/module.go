package telemetryfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/telemetry/httpserver"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"telemetry",
		fx.Provide(
			httpserver.New,
		),
	)
}
