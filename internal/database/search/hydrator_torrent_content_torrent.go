package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gen/field"
)

type torrentContentTorrentHydratorConfig struct {
	files bool
}

type HydrateTorrentContentTorrentOption func(config *torrentContentTorrentHydratorConfig)

func HydrateTorrentContentTorrentWithFiles() HydrateTorrentContentTorrentOption {
	return func(config *torrentContentTorrentHydratorConfig) {
		config.files = true
	}
}

func HydrateTorrentContentTorrent(options ...HydrateTorrentContentTorrentOption) query.Option {
	var config torrentContentTorrentHydratorConfig
	for _, option := range options {
		option(&config)
	}
	return query.HydrateHasOne[TorrentContentResultItem, model.Torrent, protocol.ID](
		torrentContentTorrentHydrator{config},
	)
}

type torrentContentTorrentHydrator struct {
	torrentContentTorrentHydratorConfig
}

func (h torrentContentTorrentHydrator) RootToSubID(root TorrentContentResultItem) (protocol.ID, bool) {
	return root.InfoHash, true
}

func (h torrentContentTorrentHydrator) GetSubs(ctx context.Context, dbCtx query.DbContext, ids []protocol.ID) ([]model.Torrent, error) {
	result, err := search{dbCtx.Query()}.Torrents(ctx, query.Where(TorrentInfoHashCriteria(ids...)), query.Preload(func(q *dao.Query) []field.RelationField {
		preload := []field.RelationField{
			q.Torrent.Sources.RelationField,
			q.Torrent.Sources.TorrentSource.RelationField,
			q.Torrent.Hint.RelationField,
			q.Torrent.Tags.RelationField,
		}
		if h.files {
			preload = append(preload, q.Torrent.Files.RelationField)
		}
		return preload
	}))
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
