package torznab

import "context"

type Client interface {
	Search(context.Context, SearchRequest) (SearchResult, error)
}
