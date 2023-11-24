package resolvers

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/persistence"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	persistence persistence.Persistence
	search      search.Search
}

func New(
	persistence persistence.Persistence,
	search search.Search,
) gql.ResolverRoot {
	return &Resolver{
		persistence: persistence,
		search:      search,
	}
}
