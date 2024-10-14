package gqlmodel

import (
	"context"
	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type QueueJobsQueryInput struct {
	Queues      []string
	Statuses    []model.QueueJobStatus
	Limit       model.NullUint
	Page        model.NullUint
	Offset      model.NullUint
	TotalCount  model.NullBool
	HasNextPage model.NullBool
	Facets      *gen.QueueJobsFacetsInput
	OrderBy     []gen.QueueJobsOrderByInput
}

type QueueJobsQueryResult struct {
	TotalCount   uint
	HasNextPage  bool
	Items        []model.QueueJob
	Aggregations gen.QueueJobsAggregations
}

func (r QueueQuery) Jobs(
	ctx context.Context,
	query QueueJobsQueryInput,
) (QueueJobsQueryResult, error) {
	limit := uint(10)
	if query.Limit.Valid {
		limit = query.Limit.Uint
	}
	options := []q.Option{
		q.SearchParams{
			Limit:       model.NullUint{Valid: true, Uint: limit},
			Page:        query.Page,
			Offset:      query.Offset,
			TotalCount:  query.TotalCount,
			HasNextPage: query.HasNextPage,
		}.Option(),
	}
	if query.Facets != nil {
		var qFacets []q.Facet
		if queue, ok := query.Facets.Queue.ValueOK(); ok {
			qFacets = append(qFacets, queueJobQueueFacet(*queue))
		}
		if status, ok := query.Facets.Status.ValueOK(); ok {
			qFacets = append(qFacets, queueJobStatusFacet(*status))
		}
		options = append(options, q.WithFacet(qFacets...))
	}
	var criteria []q.Criteria
	if query.Queues != nil {
		criteria = append(criteria, search.QueueJobQueueCriteria(query.Queues...))
	}
	if query.Statuses != nil {
		criteria = append(criteria, search.QueueJobStatusCriteria(query.Statuses...))
	}
	if len(criteria) > 0 {
		options = append(options, q.Where(criteria...))
	}
	fullOrderBy := maps.NewInsertMap[search.QueueJobsOrderBy, search.OrderDirection]()
	for _, ob := range query.OrderBy {
		direction := search.OrderDirectionAscending
		if desc, ok := ob.Descending.ValueOK(); ok && *desc {
			direction = search.OrderDirectionDescending
		}
		field, err := search.ParseQueueJobsOrderBy(ob.Field.String())
		if err != nil {
			return QueueJobsQueryResult{}, err
		}
		fullOrderBy.Set(field, direction)
	}
	options = append(options, search.QueueJobsFullOrderBy(fullOrderBy).Option())
	result, resultErr := r.QueueJobSearch.QueueJobs(ctx, options...)
	if resultErr != nil {
		return QueueJobsQueryResult{}, resultErr
	}
	return transformQueueJobsQueryResult(result)
}

func transformQueueJobsQueryResult(result q.GenericResult[model.QueueJob]) (QueueJobsQueryResult, error) {
	aggs, err := transformQueueJobsAggregations(result.Aggregations)
	if err != nil {
		return QueueJobsQueryResult{}, err
	}
	return QueueJobsQueryResult{
		TotalCount:   result.TotalCount,
		HasNextPage:  result.HasNextPage,
		Items:        result.Items,
		Aggregations: aggs,
	}, nil
}

func transformQueueJobsAggregations(aggs q.Aggregations) (gen.QueueJobsAggregations, error) {
	a := gen.QueueJobsAggregations{}
	if queue, ok := aggs[search.QueueJobQueueFacetKey]; ok {
		agg, err := queueJobQueueAggs(queue.Items)
		if err != nil {
			return a, err
		}
		a.Queue = agg
	}
	if status, ok := aggs[search.QueueJobStatusFacetKey]; ok {
		agg, err := queueJobStatusAggs(status.Items)
		if err != nil {
			return a, err
		}
		a.Status = agg
	}
	return a, nil
}
