package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
)

// QueueJobsOrderBy represents sort orders for queue jobs search results
// ENUM(created_at, ran_at, priority)
type QueueJobsOrderBy string

type QueueJobsFullOrderBy = maps.InsertMap[QueueJobsOrderBy, OrderDirection]
