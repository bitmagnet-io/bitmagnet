package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

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
