package gqlmodel

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type TorrentQuery struct {
	DaoProvider          database.DaoProvider
	TorrentMetricsClient torrentmetrics.Client
}

func (t TorrentQuery) SuggestTags(
	ctx context.Context,
	input *gen.SuggestTagsQueryInput,
) (search.TorrentSuggestTagsResult, error) {
	suggestTagsQuery := search.SuggestTagsQuery{}

	if input != nil {
		if prefix, ok := input.Prefix.ValueOK(); ok && prefix != nil {
			suggestTagsQuery.Prefix = *prefix
		}

		if exclusions, ok := input.Exclusions.ValueOK(); ok {
			suggestTagsQuery.Exclusions = exclusions
		}
	}

	return search.New(t.DaoProvider).TorrentSuggestTags(ctx, suggestTagsQuery)
}

func (t TorrentQuery) ListSources(ctx context.Context) (gen.TorrentListSourcesResult, error) {
	dao, err := t.DaoProvider.Dao()
	if err != nil {
		return gen.TorrentListSourcesResult{}, err
	}

	result, err := dao.TorrentSource.WithContext(ctx).Order(dao.TorrentSource.Key.Asc()).Find()
	if err != nil {
		return gen.TorrentListSourcesResult{}, err
	}

	sources := make([]model.TorrentSource, len(result))
	for i := range result {
		sources[i] = *result[i]
	}

	return gen.TorrentListSourcesResult{
		Sources: sources,
	}, nil
}

type TorrentMutation struct{}
