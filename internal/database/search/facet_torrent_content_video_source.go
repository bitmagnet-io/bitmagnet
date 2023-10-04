package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const VideoSourceFacetKey = "video_source"

func VideoSourceFacet(options ...query.FacetOption) query.Facet {
	return torrentContentAttributeFacet[model.VideoSource]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(VideoSourceFacetKey),
				query.FacetHasLabel("Video Source"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
		field: func(q *dao.Query) field.Field {
			return q.TorrentContent.VideoSource
		},
		parse: model.ParseVideoSource,
	}
}
