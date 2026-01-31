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
	reqOptions, reqErr := searchRequestToQueryOptions(req)
	if reqErr != nil {
		return torznab.SearchResult{}, reqErr
	}

	searchResult, searchErr := a.search.TorrentContent(
		ctx,
		search.TorrentContentDefaultOption(),
		query.WithTotalCount(false),
		query.Options(reqOptions...),
	)
	if searchErr != nil {
		return torznab.SearchResult{}, searchErr
	}

	return torrentContentResultToTorznabResult(req, searchResult), nil
}
