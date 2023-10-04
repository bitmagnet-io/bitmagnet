package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const VideoResolutionFacetKey = "video_resolution"

func videoResolutionField(q *dao.Query) field.Field {
	return q.TorrentContent.VideoResolution
}

func VideoResolutionFacet(options ...query.FacetOption) query.Facet {
	return torrentContentAttributeFacet[model.VideoResolution]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(VideoResolutionFacetKey),
				query.FacetHasLabel("Video Resolution"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
		field: videoResolutionField,
		parse: model.ParseVideoResolution,
	}
}
