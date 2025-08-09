package resolvers

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"go.uber.org/fx"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	fx.In
	Context            env.Context
	DaoProvider        database.DaoProvider
	Search             search.Search
	Workers            *registry.Registry
	Checker            health.Checker
	QueueMetricsClient queuemetrics.Client
	QueueManager       manager.Manager
	// todo: Fix this
	TorrentMetricsClient torrentmetrics.Client `optional:"true"`
	Indexer              indexer.Indexer
	PersisterAdder       persister.Adder
}
