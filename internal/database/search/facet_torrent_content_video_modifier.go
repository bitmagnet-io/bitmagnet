package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const VideoModifierFacetKey = "video_modifier"

func VideoModifierFacet(options ...query.FacetOption) query.Facet {
	return torrentContentAttributeFacet[model.VideoModifier]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(VideoModifierFacetKey),
				query.FacetHasLabel("Video Modifier"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
		field: func(q *dao.Query) field.Field {
			return q.TorrentContent.VideoModifier
		},
		parse: model.ParseVideoModifier,
	}
}
