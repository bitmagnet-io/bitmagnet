package metrics

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm"
	"strings"
	"time"
)

// BucketDuration represents the duration of bucketing for queue metrics
// ENUM(minute, hour, day)
type BucketDuration string

type Bucket struct {
	Queue           string
	Status          model.QueueJobStatus
	CreatedAtBucket time.Time
	RanAtBucket     time.Time
	Count           uint
}

type Request struct {
	BucketDuration BucketDuration
	Statuses       []model.QueueJobStatus
	Queues         []string
	CreatedFrom    time.Time
	CreatedTo      time.Time
	RanFrom        time.Time
	RanTo          time.Time
}

type Client interface {
	Request(context.Context, Request) ([]Bucket, error)
}

type client struct {
	db *gorm.DB
}

func (c client) Request(ctx context.Context, req Request) ([]Bucket, error) {
	params := []any{
		req.BucketDuration,
		req.BucketDuration,
	}
	var conditions []string
	if !req.CreatedFrom.IsZero() {
		conditions = append(conditions, "created_at >= ?")
		params = append(params, req.CreatedFrom)
	}
	if !req.CreatedTo.IsZero() {
		conditions = append(conditions, "created_at <= ?")
		params = append(params, req.CreatedTo)
	}
	if !req.RanFrom.IsZero() {
		conditions = append(conditions, "ran_at >= ?")
		params = append(params, req.RanFrom)
	}
	if !req.RanTo.IsZero() {
		conditions = append(conditions, "ran_at <= ?")
		params = append(params, req.RanTo)
	}
	if len(req.Queues) > 0 {
		conditions = append(conditions, "queue IN ?")
		params = append(params, req.Queues)
	}
	if len(req.Statuses) > 0 {
		conditions = append(conditions, "status IN ?")
		params = append(params, req.Statuses)
	}
	conditionClause := ""
	if len(conditions) > 0 {
		conditionClause = "WHERE (" + strings.Join(conditions, " AND ") + ")"
	}
	var results []Bucket
	if err := c.db.WithContext(ctx).Raw(`select queue,
       status,
       date_trunc(?, created_at) as created_at_bucket,
       date_trunc(?, ran_at) as ran_at_bucket,
       count(*) as count from queue_jobs
       `+
		conditionClause+
		`
    group by queue, status, created_at_bucket, ran_at_bucket
    order by queue, status, created_at_bucket, ran_at_bucket`,
		params...,
	).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
