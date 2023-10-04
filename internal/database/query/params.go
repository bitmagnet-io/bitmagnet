package query

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type SearchParams struct {
	QueryString model.NullString
	Limit       model.NullUint
	Offset      model.NullUint
	TotalCount  model.NullBool
	Cached      model.NullBool
}

func (s SearchParams) Option() Option {
	var options []Option
	if s.QueryString.Valid {
		options = append(options, QueryString(s.QueryString.String), OrderByQueryStringRank())
	}
	if s.Limit.Valid {
		options = append(options, Limit(s.Limit.Uint))
	}
	if s.Offset.Valid {
		options = append(options, Offset(s.Offset.Uint))
	}
	if s.TotalCount.Valid {
		options = append(options, WithTotalCount(s.TotalCount.Bool))
	}
	if s.Cached.Valid {
		if s.Cached.Bool {
			options = append(options, Cached())
		} else {
			options = append(options, CacheWarm())
		}
	}
	return Options(options...)
}
