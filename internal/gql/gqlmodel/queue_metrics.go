package gqlmodel

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
)

type QueueQuery struct {
	QueueJobSearch     search.QueueJobSearch
	QueueMetricsClient queuemetrics.Client
}

func (q QueueQuery) Metrics(ctx context.Context, input gen.QueueMetricsQueryInput) (*gen.QueueMetricsQueryResult, error) {
	req := queuemetrics.Request{}
	switch input.BucketDuration {
	case gen.MetricsBucketDurationMinute:
		req.BucketDuration = "minute"
	case gen.MetricsBucketDurationHour:
		req.BucketDuration = "hour"
	case gen.MetricsBucketDurationDay:
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
	buckets, err := q.QueueMetricsClient.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	return &gen.QueueMetricsQueryResult{
		Buckets: buckets,
	}, nil
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
