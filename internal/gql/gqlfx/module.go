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
	"github.com/bitmagnet-io/bitmagnet/internal/servarr"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
				lc servarr.Config,
				logger *zap.SugaredLogger,

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
					return resolvers.New(d, s, lc, logger), nil
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
