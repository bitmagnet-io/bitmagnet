package gqlmodel

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel/gen"
)

type TorrentQuery struct {
	TorrentSearch search.TorrentSearch
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
	return t.TorrentSearch.TorrentSuggestTags(ctx, suggestTagsQuery)
}

type TorrentMutation struct{}
