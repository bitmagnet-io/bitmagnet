package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

func ContentReleaseDateCriteria(dateRange model.DateRange) query.Criteria {
	return query.DaoCriteria{
		Conditions: func(ctx query.DBContext) ([]field.Expr, error) {
			return dateRangeConditions(ctx.Query().Content.ReleaseDate, dateRange), nil
		},
	}
}

func ContentReleaseDateCriteriaString(dateRange string) query.Criteria {
	return query.DaoCriteria{
		Conditions: func(ctx query.DBContext) ([]field.Expr, error) {
			return dateRangeConditionsStr(ctx.Query().Content.ReleaseDate, dateRange)
		},
	}
}

func dateRangeConditions(target field.Time, dateRange model.DateRange) []field.Expr {
	return []field.Expr{
		target.Gte(dateRange.StartTime()),
		target.Lt(dateRange.EndTime()),
	}
}

func dateRangeConditionsStr(target field.Time, strDateRange string) ([]field.Expr, error) {
	dateRange, dateRangeErr := model.NewDateRangeFromString(strDateRange)
	if dateRangeErr != nil {
		return nil, dateRangeErr
	}
	return dateRangeConditions(target, dateRange), nil
}
