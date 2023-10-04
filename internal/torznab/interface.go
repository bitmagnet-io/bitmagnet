package torznab

import "context"

type Client interface {
	Caps(context.Context) (Caps, error)
	Search(context.Context, SearchRequest) (SearchResult, error)
}
