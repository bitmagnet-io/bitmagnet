package adapter

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
)

func New(search search.Search) Adapter {
	return Adapter{
		search: search,
	}
}

type Adapter struct {
	search search.Search
}

func (a Adapter) Search(ctx context.Context, req torznab.SearchRequest) (torznab.SearchResult, error) {
	options := []query.Option{search.TorrentContentDefaultOption(), query.WithTotalCount(false)}
	if reqOptions, reqErr := searchRequestToQueryOptions(req); reqErr != nil {
		return torznab.SearchResult{}, reqErr
	} else {
		options = append(options, reqOptions...)
	}
	searchResult, searchErr := a.search.TorrentContent(ctx, options...)
	if searchErr != nil {
		return torznab.SearchResult{}, searchErr
	}
	return torrentContentResultToTorznabResult(req, searchResult), nil
}
