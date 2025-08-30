package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/resolvers"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/pipeline/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Schema graphql.ExecutableSchema
}

var (
	Ref = http_server.Ref.MustSub("graphql")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs the GraphQL API server"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
		builder.WithDependencies[deps](
			health.Ref,
			metrics.Ref,
			persister.Ref,
			postgres.Ref,
			processor.Ref,
			queue.Ref,
			search.Ref,
			worker.Ref,
		),
		builder.WithFxOption[deps](
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
			0,
			func(deps deps) gin.OptionFunc {
				return httpserver.New(deps.Schema)
			},
		),
	)
)
