package dhtfx

import (
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"dht",
		configfx.NewConfigModule[server.Config](server.Namespace, server.NewDefaultConfig()),
		fx.Provide(
			fx.Annotate(
				protocol.RandomNodeIDWithClientSuffix,
				fx.ResultTags(`name:"dht_node_id"`),
			),
			fx.Annotate(
				client.New,
				fx.ParamTags(`name:"dht_node_id"`),
			),
			ktable.New,
			responder.New,
			func(cfg server.Config) (socket.Adapter, error) {
				return socket.GetAdapter(cfg.SocketAdapter)
			},
			fx.Annotate(
				func(cfg server.Config, adapter socket.Adapter) socket.Runner {
					return adapter(netip.AddrPortFrom(
						netip.IPv4Unspecified(),
						cfg.Port,
					))
				},
				fx.As(new(socket.Socket)),
				fx.As(new(socket.Runner)),
			),
			fx.Annotate(
				func(sock socket.Runner) registry.Option {
					return registry.WithWorker(
						socket.Namespace,
						sock.Runner,
					)
				},
				fx.ResultTags(`group:"worker_options"`),
			),
			server.New,
			fx.Annotate(
				func(srv server.Server) registry.Option {
					return registry.WithWorker(
						server.Namespace,
						srv.Runner,
						worker.WithDependencies(socket.Namespace),
					)
				},
				fx.ResultTags(`group:"worker_options"`),
			),
			func() *concurrency.AtomicValue[server.LastResponses] {
				return &concurrency.AtomicValue[server.LastResponses]{}
			},
		),
	)
}
