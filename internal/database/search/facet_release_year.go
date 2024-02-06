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

func (r yearFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	return []query.Criteria{
		query.GenCriteria(func(ctx query.DbContext) (query.Criteria, error) {
			years := make([]uint16, 0, len(filter))
			hasNull := false
			for _, v := range filter.Values() {
				if v == "null" {
					hasNull = true
					continue
				}
				strYear := v
				vInt, intErr := strconv.Atoi(strYear)
				if intErr != nil {
					return nil, fmt.Errorf("invalid year filter specified: %w", intErr)
				} else if vInt < 1000 || vInt > 9999 {
					return nil, fmt.Errorf("out-of-bounds year filter specified: %s", strYear)
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

func (yearFacet) Values(query.FacetContext) (map[string]string, error) {
	return map[string]string{}, nil
}
