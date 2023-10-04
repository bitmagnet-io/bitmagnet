package search

import (
	"database/sql/driver"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

type torrentContentAttributeFacet[T attribute] struct {
	query.FacetConfig
	field func(*dao.Query) field.Field
	parse func(string) (T, error)
}

type attribute interface {
	fmt.Stringer
	driver.Valuer
	Label() string
}

func (f torrentContentAttributeFacet[T]) Aggregate(ctx query.FacetContext) (query.AggregationItems, error) {
	var results []struct {
		Value *T
		Count uint
	}
	q, qErr := ctx.NewAggregationQuery()
	if qErr != nil {
		return nil, qErr
	}
	fld := f.field(ctx.Query())
	if err := q.UnderlyingDB().Select(
		ctx.TableName()+"."+string(fld.ColumnName())+" as value",
		"count(*) as count",
	).Group(
		"value",
	).Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to aggregate: %w", err)
	}
	agg := make(query.AggregationItems, len(results))
	for _, item := range results {
		var key, label string
		if item.Value == nil {
			key = "null"
			label = "Unknown"
		} else {
			vV := *item.Value
			key = vV.String()
			label = vV.Label()
		}
		agg[key] = query.AggregationItem{
			Label: label,
			Count: item.Count,
		}
	}
	return agg, nil
}

func (f torrentContentAttributeFacet[T]) Criteria() []query.Criteria {
	return []query.Criteria{
		query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
			fld := f.field(ctx.Query())
			filter := f.Filter().Values()
			values := make([]driver.Valuer, 0, len(filter))
			hasNull := false
			for _, v := range filter {
				if v == "null" {
					hasNull = true
					continue
				}
				parsed, parseErr := f.parse(v)
				if parseErr != nil {
					return nil, parseErr
				}
				values = append(values, parsed)
			}
			var or []query.Criteria
			joins := maps.NewInsertMap(maps.MapEntry[string, struct{}]{Key: model.TableNameTorrentContent})
			if len(values) > 0 {
				or = append(or, query.RawCriteria{
					Query: fld.In(values...).RawExpr(),
					Joins: joins,
				})
			}
			if hasNull {
				or = append(or, query.RawCriteria{
					Query: fld.IsNull().RawExpr(),
					Joins: joins,
				})
			}
			return query.Or(or...), nil
		}),
	}
}
