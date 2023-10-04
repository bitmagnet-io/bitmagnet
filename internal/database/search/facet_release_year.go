package search

import (
	"database/sql/driver"
	"fmt"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
	"strconv"
)

const ReleaseYearFacetKey = "release_year"

func ReleaseYearFacet(options ...query.FacetOption) query.Facet {
	return yearFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(ReleaseYearFacetKey),
				query.FacetHasLabel("Release Year"),
				query.FacetUsesOrLogic(),
				// avoids counting different versions of the same piece of content, but has issues when combined with other filters:
				//query.FacetHasAggregationOption(
				//	query.Table(model.TableNameContent),
				//	ContentCoreJoins(),
				//),
			}, options...)...,
		),
		field: "release_year",
	}
}

type yearFacet struct {
	query.FacetConfig
	field string
}

func (r yearFacet) Aggregate(ctx query.FacetContext) (query.AggregationItems, error) {
	var results []struct {
		Year  string
		Count uint
	}
	q, qErr := ctx.NewAggregationQuery()
	if qErr != nil {
		return nil, qErr
	}
	if txErr := q.UnderlyingDB().Select(
		fmt.Sprintf("%s.%s as year", ctx.TableName(), r.field),
		"count(*) as count",
	).Group(
		"year",
	).Find(&results).Error; txErr != nil {
		return nil, txErr
	}
	agg := make(query.AggregationItems, len(results))
	for _, item := range results {
		key := item.Year
		label := item.Year
		if key == "" {
			key = "null"
			label = "Unknown"
		}
		agg[key] = query.AggregationItem{
			Label: label,
			Count: item.Count,
		}
	}
	return agg, nil
}

func (r yearFacet) Criteria() []query.Criteria {
	return []query.Criteria{
		query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
			filter := r.Filter().Values()
			years := make([]uint16, 0, len(filter))
			hasNull := false
			for _, v := range filter {
				if v == "null" {
					hasNull = true
					continue
				}
				strYear := v
				vInt, intErr := strconv.Atoi(strYear)
				if intErr == nil && (vInt < 1000 || vInt > 9999) {
					intErr = fmt.Errorf("out-of-bounds year filter specified: %s", strYear)
				}
				if intErr != nil {
					return nil, fmt.Errorf("invalid year filter specified: %w", intErr)
				}
				years = append(years, uint16(vInt))
			}
			yearField := ctx.Query().Content.ReleaseYear
			var or []query.Criteria
			if len(years) > 0 {
				or = append(or, query.RawCriteria{
					Query: ctx.Query().Content.UnderlyingDB().Where(yearCondition(yearField, years...).RawExpr()),
				})
			}
			if hasNull {
				or = append(or, query.RawCriteria{
					Query: ctx.Query().Content.UnderlyingDB().Where(yearField.IsNull().RawExpr()),
				})
			}
			return query.Or(or...), nil
		}),
	}
}

func yearCondition(target field.Field, years ...uint16) field.Expr {
	valuers := make([]driver.Valuer, 0, len(years))
	for _, year := range years {
		valuers = append(valuers, model.NewNullUint16(year))
	}
	return target.In(valuers...)
}
