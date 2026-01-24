package dht

import (
	"net/netip"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/socket"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Socket socket.Runner
}

var (
	Ref = ref.Root.MustSub("dht")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs a UDP socket for the DHT server"),
		builder.WithConfig[deps](Ref.MustSub("adapter"), socket.ParamAdapter()),
		builder.WithConfig[deps](Ref.MustSub("port"), socket.ParamPort),
		builder.WithFxOption[deps](fx.Provide(
			socket.GetAdapter,
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
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.Socket, nil
			},
		),
	)
)
