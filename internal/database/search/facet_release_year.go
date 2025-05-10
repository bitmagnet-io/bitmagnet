package search

import (
	"fmt"
	"strconv"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const ReleaseYearFacetKey = "release_year"

func ReleaseYearFacet(options ...query.FacetOption) query.Facet {
	return yearFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(ReleaseYearFacetKey),
				query.FacetHasLabel("Release Year"),
				query.FacetUsesOrLogic(),
				// avoids counting different versions of the same piece of content,
				// but has issues when combined with other filters:
				// query.FacetHasAggregationOption(
				// 	query.Table(model.TableNameContent),
				// 	ContentCoreJoins(),
				// ),
			}, options...)...,
		),
		field: "release_year",
	}
}

type yearFacet struct {
	query.FacetConfig
	field string
}

func (yearFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	return []query.Criteria{
		query.GenCriteria(func(ctx query.DBContext) (query.Criteria, error) {
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
			joins := maps.NewInsertMap(maps.MapEntry[string, struct{}]{Key: model.TableNameContent})
			var or []query.Criteria
			if len(years) > 0 {
				or = append(or, query.RawCriteria{
					Query: ctx.Query().
						Content.
						UnderlyingDB().
						Where(yearCondition(yearField, years...).RawExpr()),
					Joins: joins,
				})
			}
			if hasNull {
				or = append(or, query.RawCriteria{
					Query: ctx.Query().Content.UnderlyingDB().Where(yearField.IsNull().RawExpr()),
					Joins: joins,
				})
			}
			return query.Or(or...), nil
		}),
	}
}

func yearCondition(target field.Uint16, years ...uint16) field.Expr {
	return target.In(years...)
}

func (yearFacet) Values(ctx query.FacetContext) (map[string]string, error) {
	q := ctx.Query().Content

	var years []model.Year

	err := q.WithContext(ctx.Context()).Where(
		q.ReleaseYear.Gte(1000),
		q.ReleaseYear.Lte(9999),
	).Distinct(q.ReleaseYear).Pluck(q.ReleaseYear, &years)
	if err != nil {
		return nil, err
	}

	values := make(map[string]string, len(years)+1)
	values["null"] = "Unknown"

	for _, y := range years {
		values[y.String()] = y.String()
	}

	return values, nil
}
