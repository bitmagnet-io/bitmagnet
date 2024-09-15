package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gorm/clause"
)

// QueueJobsOrderBy represents sort orders for queue jobs search results
// ENUM(created_at, ran_at, priority)
type QueueJobsOrderBy string

func (ob QueueJobsOrderBy) Clauses(direction OrderDirection) []query.OrderByColumn {
	desc := direction == OrderDirectionDescending
	switch ob {
	case QueueJobsOrderByCreatedAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameQueueJob,
					Name:  "created_at",
				},
				Desc: desc,
			},
		}}
	case QueueJobsOrderByRanAt:
		return []query.OrderByColumn{{
			OrderByColumn: clause.OrderByColumn{
				Column: clause.Column{
					Table: model.TableNameQueueJob,
					Name:  "ran_at",
				},
				Desc: desc,
			},
		}}
	case QueueJobsOrderByPriority:
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

type QueueJobsFullOrderBy maps.InsertMap[QueueJobsOrderBy, OrderDirection]

func (fob QueueJobsFullOrderBy) Clauses() []query.OrderByColumn {
	im := maps.InsertMap[QueueJobsOrderBy, OrderDirection](fob)
	clauses := make([]query.OrderByColumn, 0, im.Len())
	for _, ob := range im.Entries() {
		clauses = append(clauses, ob.Key.Clauses(ob.Value)...)
	}
	return clauses
}

func (fob QueueJobsFullOrderBy) Option() query.Option {
	return query.OrderBy(fob.Clauses()...)
}
