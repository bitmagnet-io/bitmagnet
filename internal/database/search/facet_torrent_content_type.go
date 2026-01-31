package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	adapter "github.com/bitmagnet-io/bitmagnet/internal/search"
	"gorm.io/gen/field"
)

func TorrentContentTypeFacet(options ...query.FacetOption) query.Facet {
	return torrentContentTypeFacet{
		torrentContentAttributeFacet[model.ContentType]{
			FacetConfig: query.NewFacetConfig(
				append([]query.FacetOption{
					query.FacetHasKey(adapter.FacetContentType),
					query.FacetHasLabel("Content Type"),
					query.FacetUsesOrLogic(),
				}, options...)...,
			),
			field: func(q *dao.Query) field.Field {
				return field.Field(q.TorrentContent.ContentType)
			},
			parse: model.ParseContentType,
		},
	}
}

type torrentContentTypeFacet struct {
	torrentContentAttributeFacet[model.ContentType]
}

func (torrentContentTypeFacet) Values(query.FacetContext) (map[string]string, error) {
	values := make(map[string]string)
	values["null"] = "Unknown"

	for _, contentType := range model.ContentTypeValues() {
		values[string(contentType)] = contentType.Label()
	}

	return values, nil
}
