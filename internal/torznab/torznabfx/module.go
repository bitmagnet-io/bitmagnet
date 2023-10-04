package torznabfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/adapter"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/httpserver"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"torznab",
		fx.Provide(
			adapter.New,
			httpserver.New,
		),
	)
}
