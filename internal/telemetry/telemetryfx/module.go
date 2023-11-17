package telemetryfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/telemetry/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/telemetry/prometheus"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"telemetry",
		fx.Provide(
			httpserver.New,
			prometheus.New,
		),
	)
}
