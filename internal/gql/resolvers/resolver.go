package resolvers

import (
	config_manager "github.com/bitmagnet-io/bitmagnet/internal/config/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	queue_manager "github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
	QueueManager       queue_manager.Manager
	// todo: Fix this
	TorrentMetricsClient torrentmetrics.Client `optional:"true"`
	Indexer              processor.Processor
	PersisterAdder       persister.Adder
	ConfigManager        *config_manager.Manager
	Plugins              plugin.PluginInfos
	I18n                 *i18n.Bundle
}
