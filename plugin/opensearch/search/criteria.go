//go:build wasip1

package search

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/pkg/search"
	proto "github.com/bitmagnet-io/bitmagnet/proto/common/search"
	"github.com/defensestation/osquery/v2"
)

func createCriteriaFilters(params *proto.Params) ([]osquery.Mappable, error) {
	var filters []osquery.Mappable
	strCriteria := params.GetCriteria()
	if strCriteria != nil && *strCriteria != "" {
		cr, err := search.ParseCriteria([]byte(*strCriteria))
		if err != nil {
			return nil, err
		}

		filter, err := transformCriteria(cr)
		if err != nil {
			return nil, err
		}

		filters = append(filters, filter)
	}

	return filters, nil
}

func transformCriteria(criteria search.Criteria) (osquery.Mappable, error) {
	switch c := criteria.(type) {
	case search.CriteriaContentRef:
		should := make([]osquery.Mappable, 0, len(c))
		for _, cr := range c {
			should = append(should, osquery.Bool().Filter(
				osquery.Terms("contentRef.type", []any{cr.Type}),
				osquery.Terms("contentRef.source", []any{cr.Source}),
				osquery.Terms("contentRef.id", []any{cr.ID}),
			))
		}
		return osquery.Bool().Should(should...), nil
	case search.CriteriaContentType:
		return osquery.Bool().Filter(osquery.Terms("contentType", func() []any {
			var vals []any
			for _, v := range c {
				vals = append(vals, v)
			}
			return vals
		}()...)), nil
	case search.CriteriaInfoHash:
		return osquery.Bool().Filter(osquery.Terms("infoHash", func() []any {
			var vals []any
			for _, v := range c {
				vals = append(vals, v)
			}
			return vals
		}()...)), nil
	case search.And:
		var filters []osquery.Mappable
		for _, subCriteria := range c {
			filter, err := transformCriteria(subCriteria)
			if err != nil {
				return nil, err
			}
			filters = append(filters, filter)
		}
		return osquery.Bool().Must(filters...), nil
	case search.Or:
		var filters []osquery.Mappable
		for _, subCriteria := range c {
			filter, err := transformCriteria(subCriteria)
			if err != nil {
				return nil, err
			}
			filters = append(filters, filter)
		}
		return osquery.Bool().Should(filters...), nil
	case search.Not:
		var filters []osquery.Mappable
		for _, subCriteria := range c {
			filter, err := transformCriteria(subCriteria)
			if err != nil {
				return nil, err
			}
			filters = append(filters, filter)
		}
		return osquery.Bool().MustNot(filters...), nil
	default:
		return nil, errors.New("unknown criteria type")
	}
}
