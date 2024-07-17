package servarr

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/gql/gqlmodel"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gen/field"
)

type LocalSearch interface {
	TorrentContentByInfoHash(context.Context, protocol.ID) (model.TorrentContent, error)
}

type localSearch struct {
	search.Search
}

func (l localSearch) joins() query.Option {
	return query.Options(
		query.Join(func(q *dao.Query) []query.TableJoin {
			return []query.TableJoin{
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
	)
}

func (l localSearch) TorrentContentByInfoHash(ctx context.Context, infoHash protocol.ID) (gqlmodel.TorrentContent, error) {
	options := []query.Option{
		query.Where(
			search.TorrentContentInfoHashCriteria(infoHash),
		),
		l.joins(),
		search.HydrateTorrentContentContent(),
		query.Limit(1),
	}
	res, err := l.TorrentContent(ctx, options...)
	if err != nil {
		return gqlmodel.TorrentContent{}, err
	}
	return gqlmodel.NewTorrentContentFromResultItem(res.Items[0]), err
}
