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
	Latency         *time.Duration
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
	if req.Queues != nil {
		conditions = append(conditions, "queue IN ?")
		params = append(params, req.Queues)
	}
	if req.Statuses != nil {
		conditions = append(conditions, "status IN ?")
		params = append(params, req.Statuses)
	}
	conditionClause := ""
	if len(conditions) > 0 {
		conditionClause = "WHERE (" + strings.Join(conditions, " AND ") + ")"
	}
	var rawResult []rawBucket
	if err := c.db.WithContext(ctx).Raw(`select queue,
        status,
        date_trunc(?, created_at) as created_at_bucket,
        date_trunc(?, ran_at) as ran_at_bucket,
        count(*) as count,
        avg(ran_at-created_at) as latency
        from queue_jobs
       `+
		conditionClause+
		`
    group by queue, status, created_at_bucket, ran_at_bucket
    order by queue, status, created_at_bucket, ran_at_bucket`,
		params...,
	).Scan(&rawResult).Error; err != nil {
		return nil, err
	}
	result := make([]Bucket, len(rawResult))
	for i, raw := range rawResult {
		result[i] = raw.bucket()
	}
	return result, nil
}

type rawBucket struct {
	Queue           string
	Status          model.QueueJobStatus
	CreatedAtBucket time.Time
	RanAtBucket     time.Time
	Count           uint
	Latency         model.Duration
}

func (b rawBucket) bucket() Bucket {
	var latency *time.Duration
	if b.Latency > 0 {
		l := time.Duration(b.Latency)
		latency = &l
	}
	return Bucket{
		Queue:           b.Queue,
		Status:          b.Status,
		CreatedAtBucket: b.CreatedAtBucket,
		RanAtBucket:     b.RanAtBucket,
		Count:           b.Count,
		Latency:         latency,
	}
}
