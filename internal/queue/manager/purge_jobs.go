package manager

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func (m manager) PurgeJobs(ctx context.Context, req PurgeJobsRequest) error {
	if len(req.Queues) == 0 && len(req.Statuses) == 0 {
		_, err := m.db.WithContext(ctx).Raw("TRUNCATE TABLE queue_jobs;").Rows()
		return err
	}
	q := m.dao.QueueJob.WithContext(ctx)
	where := false
	if len(req.Queues) > 0 {
		q = q.Where(m.dao.QueueJob.Queue.In(req.Queues...))
		where = true
	}
	if len(req.Statuses) > 0 {
		statuses := make([]string, len(req.Statuses))
		for i, s := range req.Statuses {
			statuses[i] = string(s)
		}
		q = q.Where(m.dao.QueueJob.Status.In(statuses...))
		where = true
	}
	db := q.UnderlyingDB()
	if !where {
		db = db.Where("true")
	}
	return db.Delete(&model.QueueJob{}).Error
}
