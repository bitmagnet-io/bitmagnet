package torznabfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/adapter"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/httpserver"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"torznab",
		configfx.NewConfigModule[torznab.UntypedConfig]("torznab", torznab.NewDefaultUntypedConfig()),
		fx.Provide(
			adapter.New,
			torznab.New,
			httpserver.New,
		),
	)
}
