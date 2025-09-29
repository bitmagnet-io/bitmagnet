package search

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

type TorrentContentSearch interface {
	TorrentContent(ctx context.Context, options ...query.Option) (adapter.TorrentContentResult, error)
}

func (s search) TorrentContent(ctx context.Context, options ...query.Option) (adapter.TorrentContentResult, error) {
	return query.GenericQuery[adapter.TorrentContentResultItem](
		ctx,
		s.daoProvider,
		query.Options(append([]query.Option{query.SelectAll()}, options...)...),
		model.TableNameTorrentContent,
		func(ctx context.Context, q *dao.Query) query.SubQuery {
			return query.GenericSubQuery[dao.ITorrentContentDo]{
				SubQuery: q.TorrentContent.WithContext(ctx).ReadDB(),
			}
		},
	)
}

func TorrentContentDefaultOption() query.Option {
	return query.Options(
		query.DefaultOption(),
		TorrentContentDefaultHydrate(),
		TorrentContentCoreJoins(),
		query.OrderBy(
			query.OrderByColumn{
				OrderByColumn: clause.OrderByColumn{
					Column: clause.Column{
						Table: clause.CurrentTable,
						Name:  "published_at",
					},
					Desc: true,
				},
			},
		),
	)
}

func TorrentContentCoreJoins() query.Option {
	return query.Options(
		query.Join(func(q *dao.Query) []query.TableJoin {
			return []query.TableJoin{
				{
					Table: q.Torrent,
					On: []field.Expr{
						q.TorrentContent.InfoHash.EqCol(q.Torrent.InfoHash),
					},
					Type: query.TableJoinTypeInner,
				},
				{
					Table: q.Content,
					On: []field.Expr{
						q.TorrentContent.ContentType.EqCol(q.Content.Type),
						q.TorrentContent.ContentSource.EqCol(q.Content.Source),
						q.TorrentContent.ContentID.EqCol(q.Content.ID),
					},
					Type: query.TableJoinTypeLeft,
				},
			}
		}),
		ContentCoreJoins(),
	)
}

func TorrentContentDefaultHydrate() query.Option {
	return query.Options(
		HydrateTorrentContentTorrent(),
		HydrateTorrentContentContent(),
	)
}
