package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const TorrentContentTypeFacetKey = "content_type"

func TorrentContentTypeFacet(options ...query.FacetOption) query.Facet {
	return torrentContentAttributeFacet[model.ContentType]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(TorrentContentTypeFacetKey),
				query.FacetHasLabel("Content Type"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
		field: func(q *dao.Query) field.Field {
			return field.Field(q.TorrentContent.ContentType)
		},
		parse: model.ParseContentType,
	}
}
