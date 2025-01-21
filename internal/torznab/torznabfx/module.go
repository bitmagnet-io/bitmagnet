package torznabfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/adapter"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/settings"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"torznab",
		configfx.NewConfigModule[torznab.Config]("torznab", torznab.NewDefaultConfig()),
		fx.Provide(
			adapter.New,
			settings.New,
			httpserver.New,
		),
	)
}
