package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
	"gorm.io/gorm/clause"
)

func QueueJobsOrderByClauses(ob adapter.QueueJobsOrderBy, direction adapter.OrderDirection) []query.OrderByColumn {
	desc := direction == adapter.OrderDirectionDescending

	switch ob {
	case adapter.QueueJobsOrderByCreatedAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameQueueJob,
					Name:  "created_at",
				},
				Desc: desc,
			},
		}}
	case adapter.QueueJobsOrderByRanAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameQueueJob,
					Name:  "ran_at",
				},
				Desc: desc,
			},
		}}
	case adapter.QueueJobsOrderByPriority:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameQueueJob,
					Name:  "priority",
				},
				Desc: desc,
			},
		}}
	default:
		return []query.OrderByColumn{}
	}
}

func QueueJobsFullOrderByClauses(fob adapter.QueueJobsFullOrderBy) []query.OrderByColumn {
	im := fob
	clauses := make([]query.OrderByColumn, 0, im.Len())

	for _, ob := range im.Entries() {
		clauses = append(clauses, QueueJobsOrderByClauses(ob.Key, ob.Value)...)
	}

	return clauses
}
