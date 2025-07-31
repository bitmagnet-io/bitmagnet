package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/resolvers"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type (
	config struct{}

	deps struct {
		fx.In
		Schema graphql.ExecutableSchema
	}
)

var (
	Ref = http_server.Ref.MustSub("graphql")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithDependencies[config, deps](
			health.Ref,
			metrics.Ref,
			persister.Ref,
			postgres.Ref,
			indexer.Ref,
			queue.Ref,
			search.Ref,
			worker.Ref,
		),
		builder.WithFxOption[config, deps](
			fx.Provide(
				func(
					resolverRoot resolvers.Resolver,
				) graphql.ExecutableSchema {
					return gql.NewExecutableSchema(gql.Config{
						Resolvers: &resolverRoot,
					})
				},
			),
		),
		builder.WithGinOption(
			Ref,
			func(_ config, deps deps) gin.OptionFunc {
				return httpserver.New(deps.Schema)
			},
		),
	)
)
