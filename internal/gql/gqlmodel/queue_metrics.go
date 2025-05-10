package gqlmodel

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/queuemetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
)

type QueueQuery struct {
	QueueJobSearch     search.QueueJobSearch
	QueueMetricsClient queuemetrics.Client
}

func (qq QueueQuery) Metrics(
	ctx context.Context,
	input gen.QueueMetricsQueryInput,
) (*gen.QueueMetricsQueryResult, error) {
	req := queuemetrics.Request{
		StartTime: nilToZero(input.StartTime.Value()),
		EndTime:   nilToZero(input.EndTime.Value()),
		Statuses:  input.Statuses.Value(),
		Queues:    input.Queues.Value(),
	}

	switch input.BucketDuration {
	case gen.MetricsBucketDurationMinute:
		req.BucketDuration = "minute"
	case gen.MetricsBucketDurationHour:
		req.BucketDuration = "hour"
	case gen.MetricsBucketDurationDay:
		req.BucketDuration = "day"
	}

	buckets, err := qq.QueueMetricsClient.Request(ctx, req)
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

func (m *QueueMutation) EnqueueReprocessTorrentsBatch(
	ctx context.Context,
	input manager.EnqueueReprocessTorrentsBatchRequest,
) (*string, error) {
	err := m.QueueManager.EnqueueReprocessTorrentsBatch(ctx, input)
	return nil, err
}
