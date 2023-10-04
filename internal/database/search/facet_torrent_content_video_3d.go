package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const Video3dFacetKey = "video_3d"

func video3dField(q *dao.Query) field.Field {
	return q.TorrentContent.Video3d
}

func Video3dFacet(options ...query.FacetOption) query.Facet {
	return torrentContentAttributeFacet[model.Video3d]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(Video3dFacetKey),
				query.FacetHasLabel("Video 3D"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
		field: video3dField,
		parse: model.ParseVideo3d,
	}
}
