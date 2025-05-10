package prometheus

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// namespace used in fully-qualified metrics names.
const namespace = "bitmagnet_queue"

// queueMetricsCollector gathers queue metrics.
// It implements prometheus.Collector interface.
type queueMetricsCollector struct {
	query  lazy.Lazy[*dao.Query]
	logger *zap.SugaredLogger
}

var tasksQueuedDesc = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "", "jobs_total"),
	"Number of tasks enqueued; broken down by queue and status.",
	[]string{"queue", "status"}, nil,
)

func (qmc *queueMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(qmc, ch)
}

func (qmc *queueMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	queueInfos, err := qmc.collectQueueStatusInfos()
	if err != nil {
		qmc.logger.Errorf("Failed to collect metrics data: %s", err)
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
	q, err := qmc.query.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	var queueInfos []*queueStatusInfo

	err = q.QueueJob.WithContext(context.Background()).UnderlyingDB().Raw(
		"SELECT queue, status, count(*) FROM queue_jobs GROUP BY queue, status",
	).Find(&queueInfos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get queue status info: %w", err)
	}

	return queueInfos, nil
}
