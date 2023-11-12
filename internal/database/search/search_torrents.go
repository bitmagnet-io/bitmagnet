package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gen/field"
)

type TorrentsResult = query.GenericResult[model.Torrent]

type TorrentSearch interface {
	Torrents(ctx context.Context, options ...query.Option) (TorrentsResult, error)
	TorrentsWithMissingInfoHashes(ctx context.Context, infoHashes []protocol.ID, options ...query.Option) (TorrentsWithMissingInfoHashesResult, error)
}

func (s search) Torrents(ctx context.Context, options ...query.Option) (TorrentsResult, error) {
	return query.GenericQuery[model.Torrent](
		ctx,
		s.q,
		query.Options(options...),
		model.TableNameTorrent,
		func(ctx context.Context, q *dao.Query) query.SubQuery {
			return query.GenericSubQuery[dao.ITorrentDo]{
				SubQuery: q.Torrent.WithContext(ctx).ReadDB(),
			}
		},
	)
}

func TorrentDefaultPreload() query.Option {
	return query.Preload(func(q *dao.Query) []field.RelationField {
		return []field.RelationField{
			q.Torrent.Sources.RelationField.Order(q.TorrentsTorrentSource.CreatedAt),
			q.Torrent.Sources.TorrentSource.RelationField,
			q.Torrent.Files.RelationField.Order(q.TorrentFile.Index),
		}
	})
}

type TorrentsWithMissingInfoHashesResult struct {
	Torrents          []model.Torrent
	MissingInfoHashes []protocol.ID
}

func (s search) TorrentsWithMissingInfoHashes(ctx context.Context, infoHashes []protocol.ID, options ...query.Option) (TorrentsWithMissingInfoHashesResult, error) {
	searchResult, searchErr := s.Torrents(ctx, append([]query.Option{query.Where(TorrentInfoHashCriteria(infoHashes...))}, options...)...)
	if searchErr != nil {
		return TorrentsWithMissingInfoHashesResult{}, searchErr
	}
	torrents := make([]model.Torrent, 0, len(searchResult.Items))
	missingInfoHashes := make([]protocol.ID, 0, len(infoHashes)-len(searchResult.Items))
	foundInfoHashes := make(map[protocol.ID]struct{}, len(searchResult.Items))
nextInfoHash:
	for _, h := range infoHashes {
		for _, t := range searchResult.Items {
			if t.InfoHash == h {
				if _, ok := foundInfoHashes[h]; ok {
					continue nextInfoHash
				}
				foundInfoHashes[h] = struct{}{}
				torrents = append(torrents, t)
				continue nextInfoHash
			}
		}
		missingInfoHashes = append(missingInfoHashes, h)
	}
	return TorrentsWithMissingInfoHashesResult{
		Torrents:          torrents,
		MissingInfoHashes: missingInfoHashes,
	}, nil
}
