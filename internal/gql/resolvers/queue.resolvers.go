package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/gql"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/processor"
	"github.com/bitmagnet-io/bitmagnet/internal/queue/manager"
)

// RanAt is the resolver for the ranAt field.
func (r *queueJobResolver) RanAt(ctx context.Context, obj *model.QueueJob) (*time.Time, error) {
	if obj == nil {
		return nil, nil
	}
	t := obj.RanAt.Time
	if t.IsZero() {
		return nil, nil
	}
	return &t, nil
}

// Jobs is the resolver for the jobs field.
func (r *queueQueryResultResolver) Jobs(ctx context.Context, obj *gqlmodel.QueueQueryResult, input gqlmodel.QueueJobsQueryInput) (gqlmodel.QueueJobsQueryResult, error) {
	return gqlmodel.QueueQueryResult{QueueJobSearch: r.Search}.Jobs(ctx, input)
}

// ClassifierRematch is the resolver for the classifierRematch field.
func (r *queueEnqueueReprocessTorrentsBatchInputResolver) ClassifierRematch(ctx context.Context, obj *manager.EnqueueReprocessTorrentsBatchRequest, data *bool) error {
	if data != nil && *data {
		obj.ClassifyMode = processor.ClassifyModeRematch
	}
	return nil
}

// QueueJob returns gql.QueueJobResolver implementation.
func (r *Resolver) QueueJob() gql.QueueJobResolver { return &queueJobResolver{r} }

// QueueQueryResult returns gql.QueueQueryResultResolver implementation.
func (r *Resolver) QueueQueryResult() gql.QueueQueryResultResolver {
	return &queueQueryResultResolver{r}
}

// QueueEnqueueReprocessTorrentsBatchInput returns gql.QueueEnqueueReprocessTorrentsBatchInputResolver implementation.
func (r *Resolver) QueueEnqueueReprocessTorrentsBatchInput() gql.QueueEnqueueReprocessTorrentsBatchInputResolver {
	return &queueEnqueueReprocessTorrentsBatchInputResolver{r}
}

type queueJobResolver struct{ *Resolver }
type queueQueryResultResolver struct{ *Resolver }
type queueEnqueueReprocessTorrentsBatchInputResolver struct{ *Resolver }
