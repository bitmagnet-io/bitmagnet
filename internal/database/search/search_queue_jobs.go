package search

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
)

type QueueJobResult = adapter.Result[model.QueueJob]

type QueueJobSearch interface {
	QueueJobs(ctx context.Context, options ...query.Option) (result QueueJobResult, err error)
}

func (s search) QueueJobs(ctx context.Context, options ...query.Option) (result QueueJobResult, err error) {
	return query.GenericQuery[model.QueueJob](
		ctx,
		s.daoProvider,
		query.Options(append([]query.Option{query.SelectAll()}, options...)...),
		model.TableNameQueueJob,
		func(ctx context.Context, q *dao.Query) query.SubQuery {
			return query.GenericSubQuery[dao.IQueueJobDo]{
				SubQuery: q.QueueJob.WithContext(ctx).ReadDB(),
			}
		},
	)
}
