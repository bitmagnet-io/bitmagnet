package gqlmodel

import (
	"context"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
)

func (t TorrentQuery) Metrics(ctx context.Context, input gen.TorrentMetricsQueryInput) (*gen.TorrentMetricsQueryResult, error) {
	req := torrentmetrics.Request{}
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
	if sources, ok := input.Sources.ValueOK(); ok {
		req.Sources = sources
	}
	buckets, err := t.TorrentMetricsClient.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	return &gen.TorrentMetricsQueryResult{
		Buckets: buckets,
	}, nil
}
