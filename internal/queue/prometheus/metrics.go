package prometheus

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// namespace used in fully-qualified metrics names.
const namespace = "bitmagnet_queue"

// queueMetricsCollector gathers queue metrics.
// It implements prometheus.Collector interface.
type queueMetricsCollector struct {
	daoProvider database.DaoProvider
	logger      *zap.Logger
}

var tasksQueuedDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "", "jobs_total"),
	"Number of tasks enqueued; broken down by queue and status.",
	[]string{"queue", "status"}, nil,
)

func (*queueMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tasksQueuedDesc
}

func (qmc *queueMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	queueInfos, err := qmc.collectQueueStatusInfos()
	if err != nil {
		qmc.logger.Error("failed to collect metrics data", zap.Error(err))
	}

	for _, info := range queueInfos {
		ch <- prometheus.MustNewConstMetric(
			tasksQueuedDesc,
			prometheus.GaugeValue,
			float64(info.Count),
			info.Queue,
			info.Status.String(),
		)
	}
}

type queueStatusInfo struct {
	Queue  string
	Status model.QueueJobStatus
	Count  int
}

func (qmc *queueMetricsCollector) collectQueueStatusInfos() ([]*queueStatusInfo, error) {
	daoQ, err := qmc.daoProvider.Dao()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire database: %w", err)
	}

	var queueInfos []*queueStatusInfo

	err = daoQ.QueueJob.WithContext(context.Background()).UnderlyingDB().Raw(
		"SELECT queue, status, count(*) FROM queue_jobs GROUP BY queue, status",
	).Find(&queueInfos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get queue status info: %w", err)
	}

	return queueInfos, nil
}
