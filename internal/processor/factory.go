package processor

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocker"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/queue"
)

func New(
	defaultWorkflow classifier.Workflow,
	searchClient search.Search,
	daoProvider database.DaoTransactionProvider,
	blockerBlocker blocker.Blocker,
	runner classifier.Runner,
	queueJobProvider queue.JobProvider[MessageParams],
	persister persister.Adder,
) Processor {
	return indexer{
		searchClient:     searchClient,
		daoProvider:      daoProvider,
		blocker:          blockerBlocker,
		runner:           runner,
		defaultWorkflow:  defaultWorkflow,
		queueJobProvider: queueJobProvider,
		persister:        persister,
	}
}
