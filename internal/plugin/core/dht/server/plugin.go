package server

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht/socket"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Server server.Runner
}

var (
	Ref = dht.Ref.MustSub("server")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[deps](socket.Ref),
		builder.WithConfigParam[deps](Ref.MustSub("query_timeout"), server.ParamQueryTimeout),
		builder.WithFxOption[deps](fx.Provide(
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
			server.New,
			func(rnr server.Runner) server.Server {
				return rnr
			},
			func() *atomic.Value[server.LastResponses] {
				return &atomic.Value[server.LastResponses]{}
			},
		)),
		builder.WithWorkerRegistryOption[deps](func(deps deps) registry.Option {
			return registry.WithWorker(
				Ref.String(),
				deps.Server,
				worker.WithDependencies(socket.Ref.String()),
			)
		}),
	)
)
