package search

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/maps"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

type ContentResultItem struct {
	query.ResultItem
	model.Content
}

type ContentResult = query.GenericResult[ContentResultItem]

type ContentSearch interface {
	Content(ctx context.Context, options ...query.Option) (result ContentResult, err error)
}

func (s search) Content(ctx context.Context, options ...query.Option) (result ContentResult, err error) {
	return query.GenericQuery[ContentResultItem](
		ctx,
		s.q,
		query.Options(append([]query.Option{query.SelectAll()}, options...)...),
		model.TableNameContent,
		func(ctx context.Context, q *dao.Query) query.SubQuery {
			return query.GenericSubQuery[dao.IContentDo]{
				SubQuery: q.Content.WithContext(ctx).ReadDB(),
			}
		},
	)
}

func ContentDefaultOption() query.Option {
	return query.Options(
		query.DefaultOption(),
		ContentCoreJoins(),
		ContentDefaultPreload(),
		ContentDefaultHydrate(),
	)
}

func ContentCoreJoins() query.Option {
	return query.Join(func(q *dao.Query) []query.TableJoin {
		return []query.TableJoin{
			{
				Table: q.ContentCollectionContent,
				On: []field.Expr{
					q.ContentCollectionContent.ContentType.EqCol(q.Content.Type),
					q.ContentCollectionContent.ContentSource.EqCol(q.Content.Source),
					q.ContentCollectionContent.ContentID.EqCol(q.Content.ID),
				},
				Type: query.TableJoinTypeInner,
				Dependencies: maps.NewInsertMap(
					maps.MapEntry[string, struct{}]{Key: q.Content.TableName()},
				),
			},
			{
				Table: q.ContentCollection,
				On: []field.Expr{
					q.ContentCollection.Type.EqCol(
						q.ContentCollectionContent.ContentCollectionType,
					),
					q.ContentCollection.Source.EqCol(
						q.ContentCollectionContent.ContentCollectionSource,
					),
					q.ContentCollection.ID.EqCol(q.ContentCollectionContent.ContentCollectionID),
				},
				Type: query.TableJoinTypeInner,
				Dependencies: maps.NewInsertMap(
					maps.MapEntry[string, struct{}]{Key: q.ContentCollectionContent.TableName()},
				),
			},
			{
				Table: q.TorrentContent,
				On: []field.Expr{
					q.TorrentContent.ContentType.EqCol(q.Content.Type),
					q.TorrentContent.ContentSource.EqCol(q.Content.Source),
					q.TorrentContent.ContentID.EqCol(q.Content.ID),
				},
				Type: query.TableJoinTypeInner,
				Dependencies: maps.NewInsertMap(
					maps.MapEntry[string, struct{}]{Key: q.Content.TableName()},
				),
			},
		}
	})
}

func ContentDefaultPreload() query.Option {
	return query.Preload(func(query *dao.Query) []field.RelationField {
		return []field.RelationField{
			query.Content.MetadataSource.RelationField,
			query.Content.Attributes.RelationField,
			query.Content.Attributes.MetadataSource.RelationField,
		}
	})
}

func ContentDefaultHydrate() query.Option {
	return query.Options(
		HydrateContentCollections(),
	)
}
