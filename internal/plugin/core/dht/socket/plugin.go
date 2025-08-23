package socket

import (
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Socket socket.Runner
}

var (
	Ref = dht.Ref.MustSub("socket")

	Plugin = builder.CreatePlugin(
		Ref,
		// builder.WithDefaultConfig[deps](server.NewDefaultConfig()),
		builder.WithConfigParam[deps](Ref.MustSub("adapter"), socket.ParamAdapter),
		builder.WithConfigParam[deps](Ref.MustSub("port"), socket.ParamPort),
		builder.WithFxOption[deps](fx.Provide(
			func(adapterName socket.AdapterName) (socket.Adapter, error) {
				return socket.GetAdapter(adapterName)
			},
			fx.Annotate(
				func(adapter socket.Adapter, port socket.Port) socket.Runner {
					return adapter(netip.AddrPortFrom(
						netip.IPv4Unspecified(),
						uint16(port),
					))
				},
				fx.As(new(socket.Socket)),
				fx.As(new(socket.Runner)),
			),
		)),
		builder.WithWorkerRegistryOption(func(deps deps) registry.Option {
			return registry.WithWorker(
				Ref.String(),
				deps.Socket,
			)
		}),
	)
)
