package gqlmodel

import (
	"context"
	q "github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

type TorrentFilesQueryInput struct {
	InfoHashes  []protocol.ID
	Limit       model.NullUint
	Page        model.NullUint
	Offset      model.NullUint
	TotalCount  model.NullBool
	HasNextPage model.NullBool
	Cached      model.NullBool
	OrderBy     []gen.TorrentFilesOrderByInput
}

func (t TorrentQuery) Files(ctx context.Context, query TorrentFilesQueryInput) (search.TorrentFilesResult, error) {
	limit := uint(10)
	if query.Limit.Valid {
		limit = query.Limit.Uint
	}
	options := []q.Option{
		q.SearchParams{
			Limit:             model.NullUint{Valid: true, Uint: limit},
			Page:              query.Page,
			Offset:            query.Offset,
			TotalCount:        query.TotalCount,
			HasNextPage:       query.HasNextPage,
			AggregationBudget: model.NullFloat64{Valid: true, Float64: 0},
		}.Option(),
	}
	var criteria []q.Criteria
	if query.InfoHashes != nil {
		criteria = append(criteria, search.TorrentFileInfoHashCriteria(query.InfoHashes...))
	}
	options = append(options, q.Where(criteria...))
	fullOrderBy := maps.NewInsertMap[search.TorrentFilesOrderBy, search.OrderDirection]()
	for _, ob := range query.OrderBy {
		direction := search.OrderDirectionAscending
		if desc, ok := ob.Descending.ValueOK(); ok && *desc {
			direction = search.OrderDirectionDescending
		}
		field, err := search.ParseTorrentFilesOrderBy(ob.Field.String())
		if err != nil {
			return search.TorrentFilesResult{}, err
		}
		fullOrderBy.Set(field, direction)
	}
	options = append(options, search.TorrentFilesFullOrderBy(fullOrderBy).Option())
	return t.Search.TorrentFiles(ctx, options...)
}
