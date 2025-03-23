package search

import (
	"fmt"
	"strings"
	"gorm.io/gorm"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
)

type SizeRangeCriteria struct {
	MinBytes *int64
	MaxBytes *int64
	Key      string
}

func (c SizeRangeCriteria) Apply(q *gorm.DB) (*gorm.DB, error) {
	// If no min or max specified, return the query as is
	if c.MinBytes == nil && c.MaxBytes == nil {
		return q, nil
	}

	if c.MinBytes != nil {
		q = q.Where(fmt.Sprintf("%s >= ?", c.Key), *c.MinBytes)
	}

	if c.MaxBytes != nil {
		q = q.Where(fmt.Sprintf("%s <= ?", c.Key), *c.MaxBytes)
	}

	return q, nil
}

func (c SizeRangeCriteria) Raw(ctx query.DbContext) (query.RawCriteria, error) {
	conditions := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)

	if c.MinBytes != nil {
		conditions = append(conditions, fmt.Sprintf("%s >= ?", c.Key))
		args = append(args, *c.MinBytes)
	}

	if c.MaxBytes != nil {
		conditions = append(conditions, fmt.Sprintf("%s <= ?", c.Key))
		args = append(args, *c.MaxBytes)
	}

	return query.RawCriteria{
		Query: strings.Join(conditions, " AND "),
		Args:  args,
		Joins: maps.NewInsertMap[string, struct{}](),
	}, nil
}