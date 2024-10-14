package torrentmetrics

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Bucket struct {
	Source  string
	Bucket  time.Time
	Updated bool
	Count   uint
}

type Request struct {
	BucketDuration metrics.BucketDuration
	Sources        []string
	StartTime      time.Time
	EndTime        time.Time
	Updated        model.NullBool
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
	}
	var conditions []string
	if !req.StartTime.IsZero() {
		conditions = append(conditions, "updated_at >= ?")
		params = append(params, req.StartTime)
	}
	if !req.EndTime.IsZero() {
		conditions = append(conditions, "updated_at <= ?")
		params = append(params, req.EndTime)
	}
	if req.Sources != nil {
		conditions = append(conditions, "source IN ?")
		params = append(params, req.Sources)
	}
	if req.Updated.Valid {
		sign := ">"
		if !req.Updated.Bool {
			sign = "<="
		}
		conditions = append(conditions, "updated_at "+sign+" (created_at + interval '1 hour')")
	}
	conditionClause := ""
	if len(conditions) > 0 {
		conditionClause = "WHERE (" + strings.Join(conditions, " AND ") + ")"
	}
	var result []Bucket
	if err := c.db.WithContext(ctx).Raw(`select
        source,
        date_trunc(?, updated_at) as bucket,
        updated_at > (created_at + interval '1 hour') as updated,
        count(*) as count
        from torrents_torrent_sources
       `+
		conditionClause+
		`
    group by source, bucket, updated
    order by source, bucket, updated`,
		params...,
	).Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
