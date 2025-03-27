package manager

import (
	"context"
)

func (m manager) PurgeJobs(ctx context.Context, req PurgeJobsRequest) error {
	if len(req.Queues) == 0 && len(req.Statuses) == 0 {
		_, err := m.db.WithContext(ctx).Raw("TRUNCATE TABLE queue_jobs;").Rows()
		return err
	}

	q := m.dao.QueueJob.WithContext(ctx)
	if len(req.Queues) > 0 {
		q = q.Where(m.dao.QueueJob.Queue.In(req.Queues...))
	}

	if len(req.Statuses) > 0 {
		statuses := make([]string, len(req.Statuses))
		for i, s := range req.Statuses {
			statuses[i] = string(s)
		}

		q = q.Where(m.dao.QueueJob.Status.In(statuses...))
	}

	_, err := q.Delete()

	return err
}
