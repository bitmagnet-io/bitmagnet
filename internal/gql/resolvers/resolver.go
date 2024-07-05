package resolvers

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/servarr"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	dao           *dao.Query
	search        search.Search
	servarrConfig servarr.Config
	logger        *zap.SugaredLogger
}

func New(
	dao *dao.Query,
	search search.Search,
	servarrConfig servarr.Config,
	logger *zap.SugaredLogger,
) gql.ResolverRoot {
	return &Resolver{
		dao:           dao,
		search:        search,
		servarrConfig: servarrConfig,
		logger:        logger,
	}
}
