package gqlfx

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
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
			httpserver.New,
			func(
				ls lazy.Lazy[search.Search],
				ld lazy.Lazy[*dao.Query],
			) lazy.Lazy[gql.ResolverRoot] {
				return lazy.New(func() (gql.ResolverRoot, error) {
					s, err := ls.Get()
					if err != nil {
						return nil, err
					}
					d, err := ld.Get()
					if err != nil {
						return nil, err
					}
					return resolvers.New(d, s), nil
				})
			},
			func(
				lcfg lazy.Lazy[gql.Config],
			) lazy.Lazy[graphql.ExecutableSchema] {
				return lazy.New(func() (graphql.ExecutableSchema, error) {
					cfg, err := lcfg.Get()
					if err != nil {
						return nil, err
					}
					return gql.NewExecutableSchema(cfg), nil
				})
			},
		),
	)
}
