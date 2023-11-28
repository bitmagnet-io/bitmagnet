package resolvers

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dao    *dao.Query
	search search.Search
}

func New(
	dao *dao.Query,
	search search.Search,
) gql.ResolverRoot {
	return &Resolver{
		dao:    dao,
		search: search,
	}
}
