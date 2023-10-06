package dhtfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/crawler"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/healthcheck"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/dht/staging"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"dht",
		configfx.NewConfigModule[dht.Config]("dht_crawler", dht.NewDefaultConfig()),
		configfx.NewConfigModule[server.Config]("dht_server", server.NewDefaultConfig()),
		fx.Provide(
			crawler.New,
			healthcheck.New,
			responder.New,
			server.New,
			staging.New,
		),
	)
}
