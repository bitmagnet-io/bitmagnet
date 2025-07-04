package manager

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
)

func (m manager) PurgeJobs(ctx context.Context, req PurgeJobsRequest) error {
	db, err := m.db.GormDB()
	if err != nil {
		return err
	}

	if len(req.Queues) == 0 && len(req.Statuses) == 0 {
		_, err := db.WithContext(ctx).Raw("TRUNCATE TABLE queue_jobs;").Rows()
		return err
	}

	q := dao.Use(db).QueueJob.WithContext(ctx)
	if len(req.Queues) > 0 {
		q = q.Where(dao.QueueJob.Queue.In(req.Queues...))
	}

	if len(req.Statuses) > 0 {
		statuses := make([]string, len(req.Statuses))
		for i, s := range req.Statuses {
			statuses[i] = string(s)
		}

		q = q.Where(dao.QueueJob.Status.In(statuses...))
	}

	_, err = q.Delete()

	return err
}
