package manager

import "context"

func (m manager) PurgeJobs(ctx context.Context, req PurgeJobsRequest) error {
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
	if !where {
		q = q.Where(m.dao.QueueJob.ID.EqCol(m.dao.QueueJob.ID))
	}
	_, err := q.Delete()
	return err
}
