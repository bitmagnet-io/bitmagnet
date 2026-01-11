package graphql

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/error_registry"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/auth"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/directive"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/resolvers"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/health"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/http_server"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/postgres"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/queue"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/worker"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
	Schema         graphql.ExecutableSchema
	AuthDirectives directive.AuthDirectives
	ErrorRegistry  error_registry.Registry
	I18n           *i18n.Bundle
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
					enforcer rbac.Enforcer,
				) graphql.ExecutableSchema {
					return gql.NewExecutableSchema(gql.Config{
						Resolvers: &resolverRoot,
						Directives: gql.DirectiveRoot{
							Auth: auth.NewDirective(enforcer),
						},
					})
				},
				func(schema graphql.ExecutableSchema) *ast.Schema {
					return schema.Schema()
				},
				directive.ExtractSchemaDirectives,
				directive.ExtractAuthDirectives,
			),
		),
		builder.WithAuthObjectActions(
			func(deps deps) []rbac.ObjectAction {
				return auth.AuthObjectActions(deps.AuthDirectives)
			},
		),
		builder.WithPermissionProvider(
			func(deps) []rbac.Permission {
				return auth.Permissions()
			},
		),
		builder.WithGinOption(
			Ref,
			0,
			func(deps deps) gin.OptionFunc {
				return httpserver.New(deps.Schema, deps.ErrorRegistry, deps.I18n)
			},
		),
		builder.WithError[deps](Ref.MustSub("unauthorized"), auth.ErrUnauthorized),
	)
)
