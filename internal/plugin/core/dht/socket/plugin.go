package socket

import (
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"go.uber.org/fx"
)

type (
	Config = server.Config

	deps struct {
		fx.In
		Socket socket.Runner
	}
)

var (
	Ref = dht.Ref.MustSub("socket")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDefaultConfig[Config, deps](server.NewDefaultConfig()),
		builder.WithFxOption[Config, deps](fx.Provide(
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
		)),
		builder.WithWorkerRegistryOption(func(cfg Config, deps deps) registry.Option {
			return registry.WithWorker(
				Ref.String(),
				deps.Socket,
			)
		}),
	)
)
