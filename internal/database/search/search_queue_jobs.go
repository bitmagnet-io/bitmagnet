package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

type QueueJobResult = query.GenericResult[model.QueueJob]

type QueueJobSearch interface {
	QueueJobs(ctx context.Context, options ...query.Option) (result QueueJobResult, err error)
}

func (s search) QueueJobs(ctx context.Context, options ...query.Option) (result QueueJobResult, err error) {
	return query.GenericQuery[model.QueueJob](
		ctx,
		s.q,
		query.Options(append([]query.Option{query.SelectAll()}, options...)...),
		model.TableNameQueueJob,
		func(ctx context.Context, q *dao.Query) query.SubQuery {
			return query.GenericSubQuery[dao.IQueueJobDo]{
				SubQuery: q.QueueJob.WithContext(ctx).ReadDB(),
			}
		},
	)
}

func QueueJobQueueCriteria(queues ...string) query.Criteria {
	return query.DaoCriteria{
		Conditions: func(ctx query.DbContext) ([]field.Expr, error) {
			q := ctx.Query()
			return []field.Expr{
				q.QueueJob.Queue.In(queues...),
			}, nil
		},
	}
}

func QueueJobStatusCriteria(statuses ...model.QueueJobStatus) query.Criteria {
	strStatuses := make([]string, 0, len(statuses))
	for _, s := range statuses {
		strStatuses = append(strStatuses, s.String())
	}
	return query.DaoCriteria{
		Conditions: func(ctx query.DbContext) ([]field.Expr, error) {
			q := ctx.Query()
			return []field.Expr{
				q.QueueJob.Status.In(strStatuses...),
			}, nil
		},
	}
}
