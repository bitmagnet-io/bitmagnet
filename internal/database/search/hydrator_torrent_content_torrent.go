package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
)

func HydrateTorrentContentTorrent() query.Option {
	return query.HydrateHasOne[TorrentContentResultItem, model.Torrent, protocol.ID](
		torrentContentTorrentHydrator{},
	)
}

type torrentContentTorrentHydrator struct{}

func (h torrentContentTorrentHydrator) RootToSubID(root TorrentContentResultItem) (protocol.ID, bool) {
	return root.InfoHash, true
}

func (h torrentContentTorrentHydrator) GetSubs(ctx context.Context, dbCtx query.DbContext, ids []protocol.ID) ([]model.Torrent, error) {
	result, err := search{dbCtx.Query()}.Torrents(ctx, query.Where(TorrentInfoHashCriteria(ids...)), TorrentDefaultPreload())
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

func (h torrentContentTorrentHydrator) SubID(item model.Torrent) protocol.ID {
	return item.InfoHash
}

func (h torrentContentTorrentHydrator) Hydrate(root *TorrentContentResultItem, sub model.Torrent) {
	root.Torrent = sub
}

func (h torrentContentTorrentHydrator) MustSucceed() bool {
	return true
}
