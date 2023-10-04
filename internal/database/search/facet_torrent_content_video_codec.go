package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const VideoCodecFacetKey = "video_codec"

func VideoCodecFacet(options ...query.FacetOption) query.Facet {
	return torrentContentAttributeFacet[model.VideoCodec]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(VideoCodecFacetKey),
				query.FacetHasLabel("Video Codec"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
		field: func(q *dao.Query) field.Field {
			return q.TorrentContent.VideoCodec
		},
		parse: model.ParseVideoCodec,
	}
}
