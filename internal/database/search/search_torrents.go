package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

type TorrentsResult = query.GenericResult[model.Torrent]

type TorrentSearch interface {
	Torrents(ctx context.Context, options ...query.Option) (TorrentsResult, error)
	TorrentsWithMissingInfoHashes(ctx context.Context, infoHashes []protocol.ID, options ...query.Option) (TorrentsWithMissingInfoHashesResult, error)
	TorrentSuggestTags(ctx context.Context, query SuggestTagsQuery, options ...query.Option) (TorrentSuggestTagsResult, error)
}

func (s search) Torrents(ctx context.Context, options ...query.Option) (TorrentsResult, error) {
	return query.GenericQuery[model.Torrent](
		ctx,
		s.q,
		query.Options(append([]query.Option{query.SelectAll()}, options...)...),
		model.TableNameTorrent,
		func(ctx context.Context, q *dao.Query) query.SubQuery {
			return query.GenericSubQuery[dao.ITorrentDo]{
				SubQuery: q.Torrent.WithContext(ctx).ReadDB(),
			}
		},
	)
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

type SuggestTagsQuery struct {
	Prefix     string
	Exclusions []string
}

type SuggestedTag struct {
	Name  string
	Count int
}

type TorrentSuggestTagsResult struct {
	Suggestions []SuggestedTag
}

func (s search) TorrentSuggestTags(ctx context.Context, q SuggestTagsQuery, options ...query.Option) (TorrentSuggestTagsResult, error) {
	var criteria []query.Criteria
	if q.Prefix != "" {
		criteria = append(criteria, query.DaoCriteria{
			Conditions: func(dbCtx query.DbContext) ([]field.Expr, error) {
				return []field.Expr{
					dbCtx.Query().TorrentTag.Name.Like(q.Prefix + "%"),
				}, nil
			},
		})
	}
	if len(q.Exclusions) > 0 {
		criteria = append(criteria, query.DaoCriteria{
			Conditions: func(dbCtx query.DbContext) ([]field.Expr, error) {
				return []field.Expr{
					dbCtx.Query().TorrentTag.Name.NotIn(q.Exclusions...),
				}, nil
			},
		})
	}
	result, resultErr := query.GenericQuery[SuggestedTag](
		ctx,
		s.q,
		query.Options(append([]query.Option{
			query.Select(
				clause.Expr{
					SQL: "torrent_tags.name AS name",
				},
				clause.Expr{
					SQL: "count(torrent_tags.*) AS total_count",
				},
			),
			query.Where(criteria...),
			query.Group(
				clause.Column{
					Name: "torrent_tags.name",
				},
			),
			query.OrderBy(
				query.OrderByColumn{
					OrderByColumn: clause.OrderByColumn{
						Column: clause.Column{
							Alias: "total_count",
						},
					},
				},
				query.OrderByColumn{
					OrderByColumn: clause.OrderByColumn{
						Column: clause.Column{
							Alias: "name",
						},
					},
				},
			),
			query.Limit(10),
		}, options...)...),
		model.TableNameTorrentTag,
		func(ctx context.Context, q *dao.Query) query.SubQuery {
			return query.GenericSubQuery[dao.ITorrentTagDo]{
				SubQuery: q.TorrentTag.WithContext(ctx).ReadDB(),
			}
		},
	)
	if resultErr != nil {
		return TorrentSuggestTagsResult{}, resultErr
	}
	return TorrentSuggestTagsResult{
		Suggestions: result.Items,
	}, nil
}
