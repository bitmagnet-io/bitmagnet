package gqlfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/config"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/resolvers"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"graphql",
		fx.Provide(
			config.New,
			gql.NewExecutableSchema,
			resolvers.New,
			httpserver.New,
		),
	)
}
