package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen"
)

func HydrateContentCollections() query.Option {
	return query.HydrateHasMany[ContentResultItem, model.ContentRef, model.ContentCollectionContent, model.ContentCollection](
		contentCollectionsHydrator{},
	)
}

type contentCollectionsHydrator struct{}

func (h contentCollectionsHydrator) RootID(root ContentResultItem) (model.ContentRef, bool) {
	return model.ContentRef{
		Type:   root.Type,
		Source: root.Source,
		ID:     root.ID,
	}, true
}

func (h contentCollectionsHydrator) GetJoinSubs(ctx context.Context, dbCtx query.DBContext, ids []model.ContentRef) ([]model.ContentCollectionContent, error) {
	refMap := contentMapFromRefs(ids...)
	q := dbCtx.Query()
	var conds []gen.Condition
	for contentType, sourceMap := range refMap {
		for source, idMap := range sourceMap {
			thisIDs := make([]string, 0, len(idMap))
			for id := range idMap {
				thisIDs = append(thisIDs, id)
			}
			conds = append(conds, q.ContentCollectionContent.Where(
				q.ContentCollectionContent.ContentType.Eq(contentType),
				q.ContentCollectionContent.ContentSource.Eq(source),
				q.ContentCollectionContent.ContentID.In(thisIDs...),
			))
		}
	}
	qCtx := q.ContentCollectionContent.WithContext(ctx).Preload(
		q.ContentCollectionContent.Collection.RelationField,
		q.ContentCollectionContent.Collection.MetadataSource.RelationField,
	).UnderlyingDB()
	for _, cond := range conds {
		qCtx = qCtx.Or(cond)
	}
	var results []model.ContentCollectionContent
	if err := qCtx.Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (h contentCollectionsHydrator) JoinSubToRootIDAndSub(j model.ContentCollectionContent) (model.ContentRef, model.ContentCollection) {
	return model.ContentRef{
		Type:   j.ContentType,
		Source: j.ContentSource,
		ID:     j.ContentID,
	}, j.Collection
}

func (h contentCollectionsHydrator) Hydrate(root *ContentResultItem, subs []model.ContentCollection) {
	root.Collections = subs
}
