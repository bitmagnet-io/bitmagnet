package persistence

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type ContentPersistence interface {
	GetContent(ctx context.Context, ref model.ContentRef) (movie model.Content, err error)
}

type TorrentContentRef struct {
	model.ContentRef
	InfoHash model.Hash20
}

func (p *persistence) GetContent(ctx context.Context, ref model.ContentRef) (movie model.Content, err error) {
	m, findErr := p.q.WithContext(ctx).Content.LeftJoin(
		p.q.ContentAttribute,
		p.q.ContentAttribute.ContentType.EqCol(p.q.Content.Type),
		p.q.ContentAttribute.ContentSource.EqCol(p.q.Content.Source),
		p.q.ContentAttribute.ContentID.EqCol(p.q.Content.ID),
		p.q.ContentAttribute.Source.Eq(ref.Source),
		p.q.ContentAttribute.Key.Eq("id"),
		p.q.ContentAttribute.Value.Eq(ref.ID),
	).Where(
		p.q.Content.Where(
			p.q.Content.Type.Eq(ref.Type),
			p.q.Content.Source.Eq(ref.Source),
			p.q.Content.ID.Eq(ref.ID),
		).Or(
			p.q.ContentAttribute.Value.IsNotNull(),
		),
	).Order(
		p.q.Content.Source.Neq(ref.Source),
	).Preload(
		p.q.Content.Collections.RelationField,
		p.q.Content.Attributes.RelationField,
	).First()
	if findErr != nil {
		err = findErr
		return
	}
	movie = *m
	return
}
