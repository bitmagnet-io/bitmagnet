package processor

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
)

func New(
	config classifier.Config,
	searchClient search.Search,
	daoProvider database.DaoTransactionProvider,
	blockingManager blocking.Blocker,
	runner classifier.Runner,
) Processor {
	return processor{
		searchClient:    searchClient,
		daoProvider:     daoProvider,
		blockingManager: blockingManager,
		runner:          runner,
		defaultWorkflow: config.Workflow,
	}
}
