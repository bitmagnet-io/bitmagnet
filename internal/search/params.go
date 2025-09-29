package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Params struct {
	QueryString       model.NullString
	Limit             model.NullUint
	Page              model.NullUint
	Offset            model.NullUint
	TotalCount        model.NullBool
	HasNextPage       model.NullBool
	Cached            model.NullBool
	AggregationBudget model.NullFloat64
	OrderBy           []OrderByParam
	Criteria          Criteria
	Facets            []FacetParam
}

type OrderByParam struct {
	Key        string
	Descending bool
}
