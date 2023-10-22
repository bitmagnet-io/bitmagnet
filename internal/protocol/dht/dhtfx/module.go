package dhtfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"dht",
		configfx.NewConfigModule[server.Config]("dht_server", server.NewDefaultConfig()),
		fx.Provide(
			fx.Annotated{
				Name: "peer_id",
				Target: func() protocol.ID {
					return protocol.RandomNodeID()
				},
			},
			//healthcheck.New,
			ktable.New,
			responder.New,
			server.New,
		),
	)
}
