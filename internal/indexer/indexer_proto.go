package indexer

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/transformer/transform_model"
)

type indexerProto struct {
	indexer api.Indexer
}

func NewProto(indexer api.Indexer) Indexer {
	return indexerProto{
		indexer: indexer,
	}
}

func (i indexerProto) Add(ctx context.Context, input Input) error {
	payload := newPayload()
	input(&payload)

	_, err := i.indexer.Index(ctx, &api.IndexPayload{
		TorrentContent: slice.Map(payload.torrentContent.Values(), transform_model.TorrentContentToProto),
		DeleteTorrentContent: slice.Map(
			payload.deleteTorrentContent.Keys(),
			func(ref model.TorrentContentRef) string {
				return ref.InferID()
			},
		),
		DeleteInfoHashes: slice.Map(
			payload.deleteInfoHash.Keys(),
			func(infoHash protocol.ID) string {
				return infoHash.String()
			},
		),
	})

	return err
}
