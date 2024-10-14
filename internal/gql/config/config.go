package config

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/resolvers"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	ResolverRoot lazy.Lazy[*resolvers.Resolver]
}

func New(p Params) lazy.Lazy[gql.Config] {
	return lazy.New(func() (gql.Config, error) {
		root, err := p.ResolverRoot.Get()
		if err != nil {
			return gql.Config{}, err
		}
		return gql.Config{Resolvers: root}, nil
	})
}
