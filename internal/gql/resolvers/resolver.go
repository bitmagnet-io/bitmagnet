package resolvers

import (
	"github.com/bitmagnet-io/bitmagnet/internal/auth/api_key"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/auth/user"
	config_manager "github.com/bitmagnet-io/bitmagnet/internal/config/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/persister"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	queue_manager "github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/search/adapter/multi"
	"github.com/bitmagnet-io/bitmagnet/internal/target"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
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
	Search             multi.Search
	FacetLocalizer     search.FacetLocalizer
	Workers            *registry.Registry
	Checker            health.Checker
	QueueMetricsClient queuemetrics.Client
	QueueManager       queue_manager.Manager
	// todo: Fix this
	TorrentMetricsClient torrentmetrics.Client `optional:"true"`
	Indexer              processor.Processor
	PersisterAdder       persister.Adder
	ConfigManager        *config_manager.Manager
	Plugins              plugin.Infos
	I18n                 *i18n.Bundle
	User                 user.Service
	APIKey               api_key.Service
	RBAC                 rbac.Service
	Targets              target.Registry
}
