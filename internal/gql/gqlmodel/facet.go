package gqlmodel

import (
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/facette/natsort"
	"sort"
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

func torrentContentTypeFacet(input gen.ContentTypeFacetInput) q.Facet {
	return facet(input.Aggregate, graphql.Omittable[*model.FacetLogic]{}, input.Filter, search.TorrentContentTypeFacet)
}

func torrentSourceFacet(input gen.TorrentSourceFacetInput) q.Facet {
	var filter graphql.Omittable[[]*string]
	if f, ok := input.Filter.ValueOK(); ok {
		filterValues := make([]*string, 0, len(f))
		for _, v := range f {
			vv := v
			filterValues = append(filterValues, &vv)
		}
		filter = graphql.OmittableOf[[]*string](filterValues)
	}
	return facet(input.Aggregate, input.Logic, filter, search.TorrentSourceFacet)
}

func torrentTagFacet(input gen.TorrentTagFacetInput) q.Facet {
	var filter graphql.Omittable[[]*string]
	if f, ok := input.Filter.ValueOK(); ok {
		filterValues := make([]*string, 0, len(f))
		for _, v := range f {
			vv := v
			filterValues = append(filterValues, &vv)
		}
		filter = graphql.OmittableOf[[]*string](filterValues)
	}
	return facet(input.Aggregate, input.Logic, filter, search.TorrentTagsFacet)
}

func torrentFileTypeFacet(input gen.TorrentFileTypeFacetInput) q.Facet {
	var filter graphql.Omittable[[]*model.FileType]
	if f, ok := input.Filter.ValueOK(); ok {
		filterValues := make([]*model.FileType, 0, len(f))
		for _, v := range f {
			vv := v
			filterValues = append(filterValues, &vv)
		}
		filter = graphql.OmittableOf[[]*model.FileType](filterValues)
	}
	return facet(input.Aggregate, input.Logic, filter, search.TorrentFileTypeFacet)
}

func genreFacet(input gen.GenreFacetInput) q.Facet {
	var filter graphql.Omittable[[]*string]
	if f, ok := input.Filter.ValueOK(); ok {
		filterValues := make([]*string, 0, len(f))
		for _, v := range f {
			vv := v
			filterValues = append(filterValues, &vv)
		}
		filter = graphql.OmittableOf[[]*string](filterValues)
	}
	return facet(input.Aggregate, input.Logic, filter, search.TorrentContentGenreFacet)
}

func languageFacet(input gen.LanguageFacetInput) q.Facet {
	var filter graphql.Omittable[[]*model.Language]
	if f, ok := input.Filter.ValueOK(); ok {
		filterValues := make([]*model.Language, 0, len(f))
		for _, v := range f {
			vv := v
			filterValues = append(filterValues, &vv)
		}
		filter = graphql.OmittableOf[[]*model.Language](filterValues)
	}
	return facet(input.Aggregate, graphql.Omittable[*model.FacetLogic]{}, filter, search.TorrentContentLanguageFacet)
}

func releaseYearFacet(input gen.ReleaseYearFacetInput) q.Facet {
	return facet(input.Aggregate, graphql.Omittable[*model.FacetLogic]{}, input.Filter, search.ReleaseYearFacet)
}

func videoResolutionFacet(input gen.VideoResolutionFacetInput) q.Facet {
	return facet(input.Aggregate, graphql.Omittable[*model.FacetLogic]{}, input.Filter, search.VideoResolutionFacet)
}

func videoSourceFacet(input gen.VideoSourceFacetInput) q.Facet {
	return facet(input.Aggregate, graphql.Omittable[*model.FacetLogic]{}, input.Filter, search.VideoSourceFacet)
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

func aggs[T any, Agg comparable](
	items q.AggregationItems,
	parse func(string) (T, error),
	newAgg func(value *T, label string, count uint, isEstimate bool) Agg,
) ([]Agg, error) {
	r := make([]Agg, 0, len(items))
	labelMap := make(map[Agg]string, len(items))
	for key, item := range items {
		if key != "null" {
			v, err := parse(key)
			if err != nil {
				return nil, fmt.Errorf("error parsing aggregation item: %w", err)
			}
			agg := newAgg(&v, item.Label, item.Count, item.IsEstimate)
			r = append(r, agg)
			labelMap[agg] = item.Label
		}
	}
	sort.Slice(r, func(i, j int) bool {
		return natsort.Compare(labelMap[r[i]], labelMap[r[j]])
	})
	if null, nullOk := items["null"]; nullOk {
		r = append(r, newAgg(nil, null.Label, null.Count, null.IsEstimate))
	}
	return r, nil
}

func contentTypeAggs(items q.AggregationItems) ([]gen.ContentTypeAgg, error) {
	return aggs(items, model.ParseContentType, func(value *model.ContentType, label string, count uint, isEstimate bool) gen.ContentTypeAgg {
		return gen.ContentTypeAgg{Value: value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func torrentSourceAggs(items q.AggregationItems) ([]gen.TorrentSourceAgg, error) {
	return aggs(items, func(s string) (string, error) { return s, nil }, func(value *string, label string, count uint, isEstimate bool) gen.TorrentSourceAgg {
		return gen.TorrentSourceAgg{Value: *value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func torrentTagAggs(items q.AggregationItems) ([]gen.TorrentTagAgg, error) {
	return aggs(items, func(s string) (string, error) { return s, nil }, func(value *string, label string, count uint, isEstimate bool) gen.TorrentTagAgg {
		return gen.TorrentTagAgg{Value: *value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func torrentFileTypeAggs(items q.AggregationItems) ([]gen.TorrentFileTypeAgg, error) {
	return aggs(items, model.ParseFileType, func(value *model.FileType, label string, count uint, isEstimate bool) gen.TorrentFileTypeAgg {
		return gen.TorrentFileTypeAgg{Value: *value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func languageAggs(items q.AggregationItems) ([]gen.LanguageAgg, error) {
	return aggs(items, func(str string) (model.Language, error) {
		lang := model.ParseLanguage(str)
		if !lang.Valid {
			return "", errors.New("invalid language")
		}
		return lang.Language, nil
	}, func(value *model.Language, label string, count uint, isEstimate bool) gen.LanguageAgg {
		return gen.LanguageAgg{Value: *value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func genreAggs(items q.AggregationItems) ([]gen.GenreAgg, error) {
	return aggs(items, func(s string) (string, error) { return s, nil }, func(value *string, label string, count uint, isEstimate bool) gen.GenreAgg {
		return gen.GenreAgg{Value: *value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func releaseYearAggs(items q.AggregationItems) ([]gen.ReleaseYearAgg, error) {
	return aggs(items, model.ParseYear, func(value *model.Year, label string, count uint, isEstimate bool) gen.ReleaseYearAgg {
		return gen.ReleaseYearAgg{Value: value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func videoResolutionAggs(items q.AggregationItems) ([]gen.VideoResolutionAgg, error) {
	return aggs(items, model.ParseVideoResolution, func(value *model.VideoResolution, label string, count uint, isEstimate bool) gen.VideoResolutionAgg {
		return gen.VideoResolutionAgg{Value: value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func videoSourceAggs(items q.AggregationItems) ([]gen.VideoSourceAgg, error) {
	return aggs(items, model.ParseVideoSource, func(value *model.VideoSource, label string, count uint, isEstimate bool) gen.VideoSourceAgg {
		return gen.VideoSourceAgg{Value: value, Label: label, Count: int(count), IsEstimate: isEstimate}
	})
}

func queueJobQueueAggs(items q.AggregationItems) ([]gen.QueueJobQueueAgg, error) {
	return aggs(items, func(s string) (string, error) { return s, nil }, func(value *string, label string, count uint, isEstimate bool) gen.QueueJobQueueAgg {
		return gen.QueueJobQueueAgg{Value: *value, Label: label, Count: int(count)}
	})
}

func queueJobStatusAggs(items q.AggregationItems) ([]gen.QueueJobStatusAgg, error) {
	return aggs(items, model.ParseQueueJobStatus, func(value *model.QueueJobStatus, label string, count uint, isEstimate bool) gen.QueueJobStatusAgg {
		return gen.QueueJobStatusAgg{Value: *value, Label: label, Count: int(count)}
	})
}
