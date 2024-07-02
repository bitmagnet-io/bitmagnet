package gqlmodel

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
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
	if t, ok := input.CreatedFrom.ValueOK(); ok && !t.IsZero() {
		req.CreatedFrom = *t
	}
	if t, ok := input.CreatedTo.ValueOK(); ok && !t.IsZero() {
		req.CreatedTo = *t
	}
	if t, ok := input.RanFrom.ValueOK(); ok && !t.IsZero() {
		req.RanFrom = *t
	}
	if t, ok := input.RanTo.ValueOK(); ok && !t.IsZero() {
		req.RanTo = *t
	}
	if statuses, ok := input.Statuses.ValueOK(); ok {
		req.Statuses = statuses
	}
	if queues, ok := input.Queues.ValueOK(); ok {
		req.Queues = queues
	}
	return q.QueueMetricsClient.Request(ctx, req)
}
