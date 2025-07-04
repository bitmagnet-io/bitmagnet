package prometheus

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/logging"
	"github.com/prometheus/client_golang/prometheus"
)

func New(
	daoProvider database.DaoProvider,
	logger logging.Logger,
) prometheus.Collector {
	return &queueMetricsCollector{
		daoProvider: daoProvider,
		logger:      logger,
	}
}
