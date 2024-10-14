package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"gorm.io/gen/field"
)

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
