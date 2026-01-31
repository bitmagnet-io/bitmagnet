package target

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
)

type Registry struct {
	targets ref.Map[TorrentContentTarget]
	search  search.TorrentContent
}

func NewRegistry(targets []TorrentContentTarget, search search.TorrentContent) Registry {
	targetMap := ref.NewMap[TorrentContentTarget]()
	for _, target := range targets {
		targetMap.Set(target.Ref(), target)
	}

	return Registry{
		targets: targetMap,
		search:  search,
	}
}

type Params struct {
	Target     ref.Ref
	Index      ref.Nullable
	InfoHashes []protocol.ID
	Data       json_schema.JSONValue
}

type Result struct {
	Data              *json_schema.JSONValue
	MissingInfoHashes []protocol.ID
}

func (r Registry) Targets() []TorrentContentTarget {
	return r.targets.Values()
}

func (r Registry) Send(ctx context.Context, params Params) (Result, error) {
	target, ok := r.targets.GetOK(params.Target)
	if !ok {
		return Result{}, fmt.Errorf("%w: %w: %s", Err, ErrUnknownTarget, params.Target)
	}

	result, err := r.search.TorrentContent(ctx, search.Params{
		Index:    params.Index,
		Criteria: search.CriteriaInfoHash(params.InfoHashes),
	})
	if err != nil {
		return Result{}, fmt.Errorf("%w %s: %w: %w", Err, params.Target, ErrLookupTorrents, err)
	}

	var (
		torrents          []model.TorrentContent
		missingInfoHashes []protocol.ID
	)

	for _, infoHash := range params.InfoHashes {
		var torrent *model.TorrentContent

		for _, item := range result.Items {
			if item.InfoHash == infoHash {
				torrent = &item.TorrentContent
				break
			}
		}

		if torrent != nil {
			torrents = append(torrents, *torrent)
		} else {
			missingInfoHashes = append(missingInfoHashes, infoHash)
		}
	}

	if len(torrents) == 0 {
		return Result{}, fmt.Errorf("%w %s: %w", Err, params.Target, ErrNoTorrents)
	}

	data, err := target.Send(ctx, torrents, params.Data)
	if err != nil {
		return Result{}, fmt.Errorf("%w %s: %w: %w", Err, params.Target, ErrSend, err)
	}

	return Result{
		Data:              data,
		MissingInfoHashes: missingInfoHashes,
	}, nil
}
