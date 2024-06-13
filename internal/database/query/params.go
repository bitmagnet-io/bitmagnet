package query

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type SearchParams struct {
	QueryString       model.NullString
	Limit             model.NullUint
	Page              model.NullUint
	Offset            model.NullUint
	TotalCount        model.NullBool
	HasNextPage       model.NullBool
	Cached            model.NullBool
	AggregationBudget model.NullFloat64
}

func (s SearchParams) Option() Option {
	var options []Option
	if s.QueryString.Valid {
		options = append(options, QueryString(s.QueryString.String), OrderByQueryStringRank())
	}
	offset := uint(0)
	if s.Limit.Valid {
		options = append(options, Limit(s.Limit.Uint))
		if s.Page.Valid && s.Page.Uint > 0 {
			offset += (s.Page.Uint - 1) * s.Limit.Uint
		}
	}
	if s.Offset.Valid {
		offset += s.Offset.Uint
	}
	if offset > 0 {
		options = append(options, Offset(offset))
	}
	if s.TotalCount.Valid {
		options = append(options, WithTotalCount(s.TotalCount.Bool))
	}
	if s.HasNextPage.Valid {
		options = append(options, WithHasNextPage(s.HasNextPage.Bool))
	}
	if s.Cached.Valid {
		if s.Cached.Bool {
			options = append(options, Cached())
		} else {
			options = append(options, CacheWarm())
		}
	}
	if s.AggregationBudget.Valid {
		options = append(options, WithAggregationBudget(s.AggregationBudget.Float64))
	}
	return Options(options...)
}
