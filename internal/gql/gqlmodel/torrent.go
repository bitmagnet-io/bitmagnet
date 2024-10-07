package gqlmodel

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
)

type TorrentQuery struct {
	Search               search.Search
	TorrentMetricsClient torrentmetrics.Client
}

func (t TorrentQuery) SuggestTags(ctx context.Context, query *gen.SuggestTagsQueryInput) (search.TorrentSuggestTagsResult, error) {
	suggestTagsQuery := search.SuggestTagsQuery{}
	if query != nil {
		if prefix, ok := query.Prefix.ValueOK(); ok && prefix != nil {
			suggestTagsQuery.Prefix = *prefix
		}
		if exclusions, ok := query.Exclusions.ValueOK(); ok {
			suggestTagsQuery.Exclusions = exclusions
		}
	}
	return t.Search.TorrentSuggestTags(ctx, suggestTagsQuery)
}

type TorrentMutation struct{}
