package gqlfx

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/config"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/httpserver"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/resolvers"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"graphql",
		fx.Provide(
			config.New,
			httpserver.New,
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
		fx.Provide(
			func(p Params) Result {
				return Result{
					Resolver: lazy.New(func() (*resolvers.Resolver, error) {
						ch, err := p.Checker.Get()
						if err != nil {
							return nil, err
						}
						s, err := p.Search.Get()
						if err != nil {
							return nil, err
						}
						d, err := p.Dao.Get()
						if err != nil {
							return nil, err
						}
						return &resolvers.Resolver{
							Dao:     d,
							Search:  s,
							Checker: ch,
						}, nil
					}),
				}
			},
		),
		// inject resolver dependencies avoiding a circular dependency:
		fx.Invoke(func(
			resolver lazy.Lazy[*resolvers.Resolver],
			workers worker.Registry,
		) {
			resolver.Decorate(func(r *resolvers.Resolver) (*resolvers.Resolver, error) {
				r.Workers = workers
				return r, nil
			})
		}),
	)
}

type Params struct {
	fx.In
	Search  lazy.Lazy[search.Search]
	Dao     lazy.Lazy[*dao.Query]
	Checker lazy.Lazy[health.Checker]
}

type Result struct {
	fx.Out
	Resolver lazy.Lazy[*resolvers.Resolver]
}
