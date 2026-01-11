package transform_to_proto

import (
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	internal_search "github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	proto_search "github.com/bitmagnet-io/bitmagnet/proto/common/search"
)

func SearchParams(params internal_search.Params) *proto_search.Params {
	var (
		queryString       *string
		criteria          *string
		limit             *int32
		page              *int32
		offset            *int32
		totalCount        *bool
		hasNextPage       *bool
		aggregationBudget *float32
	)
	if params.QueryString.Valid {
		queryString = &params.QueryString.String
	}
	if bytes, err := json.Marshal(params.Criteria); err == nil {
		str := string(bytes)
		criteria = &str
	}
	if params.Limit.Valid {
		l := int32(params.Limit.Uint)
		limit = &l
	}
	if params.Page.Valid {
		p := int32(params.Page.Uint)
		page = &p
	}
	if params.Offset.Valid {
		o := int32(params.Offset.Uint)
		offset = &o
	}
	if params.TotalCount.Valid {
		totalCount = &params.TotalCount.Bool
	}
	if params.HasNextPage.Valid {
		hasNextPage = &params.HasNextPage.Bool
	}
	if params.AggregationBudget.Valid {
		ab := float32(params.AggregationBudget.Float64)
		aggregationBudget = &ab
	}
	return &proto_search.Params{
		QueryString:       queryString,
		Criteria:          criteria,
		Limit:             limit,
		Page:              page,
		Offset:            offset,
		TotalCount:        totalCount,
		HasNextPage:       hasNextPage,
		AggregationBudget: aggregationBudget,
		OrderBy:           slice.Map(params.OrderBy, OrderByParam),
		Facets:            slice.Map(params.Facets, FacetParam),
	}
}

func OrderByParam(param internal_search.OrderByParam) *proto_search.OrderByParam {
	return &proto_search.OrderByParam{
		Key:        param.Key,
		Descending: &param.Descending,
	}
}

func FacetParam(param internal_search.FacetParam) *proto_search.FacetParam {
	return &proto_search.FacetParam{
		Key:       string(param.Key),
		Filter:    param.Filter,
		Aggregate: &param.Aggregate,
		Logic:     transformNullFacetLogic(param.Logic),
	}
}

func transformNullFacetLogic(logic model.NullFacetLogic) *proto_search.FacetLogic {
	switch {
	case !logic.Valid:
		return nil
	case logic.FacetLogic == model.FacetLogicOr:
		l := proto_search.FacetLogic_or
		return &l
	case logic.FacetLogic == model.FacetLogicAnd:
		l := proto_search.FacetLogic_and
		return &l
	default:
		return nil
	}
}
