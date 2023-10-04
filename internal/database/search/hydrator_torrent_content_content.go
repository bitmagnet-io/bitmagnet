package search

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

func HydrateTorrentContentContent() query.Option {
	return query.HydrateHasOne[TorrentContentResultItem, model.Content, model.ContentRef](
		torrentContentContentHydrator{},
	)
}

type torrentContentContentHydrator struct{}

func (h torrentContentContentHydrator) RootToSubID(root TorrentContentResultItem) (model.ContentRef, bool) {
	ref := root.EntityReference()
	return ref.Val, ref.Valid
}

func (h torrentContentContentHydrator) GetSubs(ctx context.Context, dbCtx query.DbContext, ids []model.ContentRef) ([]model.Content, error) {
	contentResult, contentErr := search{dbCtx.Query()}.Content(
		ctx,
		query.Where(ContentCanonicalIdentifierCriteria(ids...)),
		query.Preload(func(query *dao.Query) []field.RelationField {
			return []field.RelationField{
				query.Content.MetadataSource.RelationField,
				query.Content.Attributes.RelationField,
				query.Content.Attributes.MetadataSource.RelationField,
			}
		}),
		ContentDefaultHydrate(),
	)
	if contentErr != nil {
		return nil, contentErr
	}
	content := make([]model.Content, 0, len(contentResult.Items))
	for _, c := range contentResult.Items {
		content = append(content, c.Content)
	}
	return content, nil
}

func (h torrentContentContentHydrator) SubID(item model.Content) model.ContentRef {
	return item.Ref()
}

func (h torrentContentContentHydrator) Hydrate(root *TorrentContentResultItem, sub model.Content) {
	root.Content = sub
}

func (h torrentContentContentHydrator) MustSucceed() bool {
	return true
}
