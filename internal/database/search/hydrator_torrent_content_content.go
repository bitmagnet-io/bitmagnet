package search

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

func HydrateTorrentContentContent() query.Option {
	return query.HydrateHasOne[TorrentContentResultItem, model.Content, model.ContentRef](
		torrentContentContentHydrator{},
	)
}

type torrentContentContentHydrator struct{}

func (torrentContentContentHydrator) RootToSubID(root TorrentContentResultItem) (model.ContentRef, bool) {
	ref := root.ContentRef()
	return ref.Val, ref.Valid
}

func (torrentContentContentHydrator) GetSubs(
	ctx context.Context,
	dbCtx query.DBContext,
	ids []model.ContentRef,
) ([]model.Content, error) {
	contentResult, contentErr := search{dbCtx.Query()}.Content(
		ctx,
		query.Where(ContentCanonicalIdentifierCriteria(ids...)),
		ContentDefaultPreload(),
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

func (torrentContentContentHydrator) SubID(item model.Content) model.ContentRef {
	return item.Ref()
}

func (torrentContentContentHydrator) Hydrate(root *TorrentContentResultItem, sub model.Content) {
	root.Content = sub
}

func (torrentContentContentHydrator) MustSucceed() bool {
	return true
}
