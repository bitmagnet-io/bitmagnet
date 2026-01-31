package proto

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/transform_from_proto"
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/transform_to_proto"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
)

type adapterProto struct {
	search.Base
	adapter api.SearchAdapter
}

func New(adapter api.SearchAdapter) search.Base {
	return &adapterProto{
		adapter: adapter,
	}
}

func (a *adapterProto) TorrentContent(
	ctx context.Context,
	params search.Params,
) (search.TorrentContentResult, error) {
	result, err := a.adapter.SearchTorrentContent(ctx, transform_to_proto.SearchParams(params))
	if err != nil {
		return search.TorrentContentResult{}, err
	}

	return transform_from_proto.SearchTorrentContentResult(result), nil
}

func (a *adapterProto) TorrentFiles(
	ctx context.Context,
	params search.Params,
) (search.TorrentFilesResult, error) {
	result, err := a.adapter.SearchTorrentFiles(ctx, transform_to_proto.SearchParams(params))
	if err != nil {
		return search.TorrentFilesResult{}, err
	}

	return transform_from_proto.SearchTorrentFilesResult(result), nil
}
