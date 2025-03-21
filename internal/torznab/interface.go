package torznab

import "context"

type Client interface {
	Caps(context.Context, Profile) (Caps, error)
	Search(context.Context, SearchRequest) (SearchResult, error)
}
