package prometheus

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func New(
	daoProvider database.DaoProvider,
	logger *zap.Logger,
) prometheus.Collector {
	return &queueMetricsCollector{
		daoProvider: daoProvider,
		logger:      logger,
	}
}
