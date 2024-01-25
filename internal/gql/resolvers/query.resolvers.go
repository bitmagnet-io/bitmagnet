package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
)

// Torrent is the resolver for the torrent field.
func (r *queryResolver) Torrent(ctx context.Context) (gqlmodel.TorrentQuery, error) {
	return gqlmodel.TorrentQuery{
		TorrentSearch: r.search,
	}, nil
}

// TorrentContent is the resolver for the torrentContent field.
func (r *queryResolver) TorrentContent(ctx context.Context) (gqlmodel.TorrentContentQuery, error) {
	return gqlmodel.TorrentContentQuery{
		TorrentContentSearch: r.search,
	}, nil
}

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
