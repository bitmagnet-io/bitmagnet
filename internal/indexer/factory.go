package indexer

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
)

func New(
	config classifier.Config,
	searchClient search.Search,
	daoProvider database.DaoTransactionProvider,
	blockerBlocker blocker.Blocker,
	runner classifier.Runner,
	queueJobProvider queue.JobProvider[MessageParams],
	persister persister.Adder,
) Processor {
	return processor{
		searchClient:     searchClient,
		daoProvider:      daoProvider,
		blocker:          blockerBlocker,
		runner:           runner,
		defaultWorkflow:  config.Workflow,
		queueJobProvider: queueJobProvider,
		persister:        persister,
	}
}
