package persister

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/indexer"
)

type flusherIndexer struct {
	iflusher
	indexer indexer.Indexer
}

func (f *flusherIndexer) flush(ctx context.Context, payload *payload) (AllTablesStats, error) {
	stats, err := f.iflusher.flush(ctx, payload)
	if err != nil {
		return stats, err
	}

	var inputs indexer.Inputs

	if payload.torrentContents.Len() > 0 {
		inputs = append(inputs, indexer.InputTorrentContent(payload.torrentContents.Values()...))
	}

	if payload.deleteInfoHashes.Len() > 0 {
		inputs = append(inputs, indexer.InputDeleteInfoHashes(payload.deleteInfoHashes.Keys()...))
	}

	if payload.deleteTorrentContent.Len() > 0 {
		inputs = append(inputs, indexer.InputDeleteTorrentContent(payload.deleteTorrentContent.Keys()...))
	}

	if len(inputs) > 0 {
		err = f.indexer.Add(ctx, inputs.Input())
	}

	return stats, err
}
