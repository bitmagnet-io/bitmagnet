package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

// Torrent is the resolver for the torrent field.
func (r *mutationResolver) Torrent(ctx context.Context) (gqlmodel.TorrentMutation, error) {
	return gqlmodel.TorrentMutation{}, nil
}

// Queue is the resolver for the queue field.
func (r *mutationResolver) Queue(ctx context.Context) (gqlmodel.QueueMutation, error) {
	return gqlmodel.QueueMutation{QueueManager: r.QueueManager}, nil
}

// Delete is the resolver for the delete field.
func (r *torrentMutationResolver) Delete(ctx context.Context, obj *gqlmodel.TorrentMutation, infoHashes []protocol.ID) (*string, error) {
	_, err := r.Dao.DeleteAndBlockTorrents(ctx, infoHashes)
	return nil, err
}

// PutTags is the resolver for the putTags field.
func (r *torrentMutationResolver) PutTags(ctx context.Context, obj *gqlmodel.TorrentMutation, infoHashes []protocol.ID, tagNames []string) (*string, error) {
	return nil, r.Dao.TorrentTag.Put(ctx, infoHashes, tagNames)
}

// SetTags is the resolver for the setTags field.
func (r *torrentMutationResolver) SetTags(ctx context.Context, obj *gqlmodel.TorrentMutation, infoHashes []protocol.ID, tagNames []string) (*string, error) {
	return nil, r.Dao.TorrentTag.Set(ctx, infoHashes, tagNames)
}

// DeleteTags is the resolver for the deleteTags field.
func (r *torrentMutationResolver) DeleteTags(ctx context.Context, obj *gqlmodel.TorrentMutation, infoHashes []protocol.ID, tagNames []string) (*string, error) {
	return nil, r.Dao.TorrentTag.Delete(ctx, infoHashes, tagNames)
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

// TorrentMutation returns gql.TorrentMutationResolver implementation.
func (r *Resolver) TorrentMutation() gql.TorrentMutationResolver { return &torrentMutationResolver{r} }

type mutationResolver struct{ *Resolver }
type torrentMutationResolver struct{ *Resolver }
