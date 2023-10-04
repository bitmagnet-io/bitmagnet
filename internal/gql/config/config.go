package config

import (
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	ResolverRoot gql.ResolverRoot
}

func New(p Params) gql.Config {
	return gql.Config{Resolvers: p.ResolverRoot}
}
