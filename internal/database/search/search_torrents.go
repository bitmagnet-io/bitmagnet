package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

type TorrentsResult = query.GenericResult[model.Torrent]

type TorrentSearch interface {
	Torrents(ctx context.Context, options ...query.Option) (result TorrentsResult, err error)
}

func (s search) Torrents(ctx context.Context, options ...query.Option) (result TorrentsResult, err error) {
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
