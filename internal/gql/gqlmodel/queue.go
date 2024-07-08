package gqlmodel

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/metrics"
)

type QueueQueryResult struct {
	QueueMetricsClient metrics.Client
}

func (q QueueQueryResult) Metrics(ctx context.Context, input gen.QueueMetricsQueryInput) ([]metrics.Bucket, error) {
	req := metrics.Request{}
	switch input.BucketDuration {
	case gen.QueueMetricsBucketDurationMinute:
		req.BucketDuration = "minute"
	case gen.QueueMetricsBucketDurationHour:
		req.BucketDuration = "hour"
	case gen.QueueMetricsBucketDurationDay:
		req.BucketDuration = "day"
	default:
		return nil, fmt.Errorf("invalid bucket duration: %s", input.BucketDuration)
	}
	if t, ok := input.StartTime.ValueOK(); ok && !t.IsZero() {
		req.StartTime = *t
	}
	if t, ok := input.EndTime.ValueOK(); ok && !t.IsZero() {
		req.EndTime = *t
	}
	if statuses, ok := input.Statuses.ValueOK(); ok {
		req.Statuses = statuses
	}
	if queues, ok := input.Queues.ValueOK(); ok {
		req.Queues = queues
	}
	return q.QueueMetricsClient.Request(ctx, req)
}

type QueueMutation struct {
	QueueManager manager.Manager
}

func (m *QueueMutation) PurgeJobs(ctx context.Context, input manager.PurgeJobsRequest) (*string, error) {
	err := m.QueueManager.PurgeJobs(ctx, input)
	return nil, err
}

func (m *QueueMutation) EnqueueReprocessTorrentsBatch(ctx context.Context, input manager.EnqueueReprocessTorrentsBatchRequest) (*string, error) {
	err := m.QueueManager.EnqueueReprocessTorrentsBatch(ctx, input)
	return nil, err
}
