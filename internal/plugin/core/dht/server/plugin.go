package server

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/dht"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/client"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/responder"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/server"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Server server.Runner
}

var (
	Ref = dht.Ref.MustSub("server")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs a DHT server node"),
		builder.WithDependencies[deps](dht.Ref),
		builder.WithConfig[deps](Ref.MustSub("query_timeout"), server.ParamQueryTimeout),
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
		builder.WithWorker(
			func(deps deps) (runner.Provider, worker.Option) {
				return deps.Server, worker.WithDependencies(dht.Ref)
			},
		),
	)
)
