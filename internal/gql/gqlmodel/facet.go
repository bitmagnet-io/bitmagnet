package gqlmodel

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func facet[T comparable](
	aggregate graphql.Omittable[*bool],
	logic graphql.Omittable[*model.FacetLogic],
	filter graphql.Omittable[[]*T],
	fn func(...q.FacetOption) q.Facet,
) q.Facet {
	var facetOptions []q.FacetOption
	if agg, aggOk := aggregate.ValueOK(); aggOk && *agg {
		facetOptions = append(facetOptions, q.FacetIsAggregated())
	}

	if l, logicOk := logic.ValueOK(); logicOk {
		facetOptions = append(facetOptions, q.FacetUsesLogic(*l))
	}

	if flt, filterOk := filter.ValueOK(); filterOk {
		f := make(q.FacetFilter, len(flt))

		for _, v := range flt {
			if v == nil {
				f["null"] = struct{}{}
			} else {
				vv := *v
				f[toString(vv)] = struct{}{}
			}
		}

		facetOptions = append(facetOptions, q.FacetHasFilter(f))
	}

	return fn(facetOptions...)
}

func toString(v any) string {
	if pStr, ok := v.(*string); ok {
		return *pStr
	}

	if str, ok := v.(string); ok {
		return str
	}

	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String()
	}

	return fmt.Sprintf("%v", v)
}

func queueJobQueueFacet(input gen.QueueJobQueueFacetInput) q.Facet {
	var filter graphql.Omittable[[]*string]

	if f, ok := input.Filter.ValueOK(); ok {
		filterValues := make([]*string, 0, len(f))

		for _, v := range f {
			vv := v
			filterValues = append(filterValues, &vv)
		}

		filter = graphql.OmittableOf[[]*string](filterValues)
	}

	return facet(input.Aggregate, graphql.Omittable[*model.FacetLogic]{}, filter, search.QueueJobQueueFacet)
}

func queueJobStatusFacet(input gen.QueueJobStatusFacetInput) q.Facet {
	var filter graphql.Omittable[[]*model.QueueJobStatus]

	if f, ok := input.Filter.ValueOK(); ok {
		filterValues := make([]*model.QueueJobStatus, 0, len(f))

		for _, v := range f {
			vv := v
			filterValues = append(filterValues, &vv)
		}

		filter = graphql.OmittableOf[[]*model.QueueJobStatus](filterValues)
	}

	return facet(input.Aggregate, graphql.Omittable[*model.FacetLogic]{}, filter, search.QueueJobStatusFacet)
}
