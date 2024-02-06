package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const VideoSourceFacetKey = "video_source"

func VideoSourceFacet(options ...query.FacetOption) query.Facet {
	return videoSourceFacet{torrentContentAttributeFacet[model.VideoSource]{
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
	}}
}

type videoSourceFacet struct {
	torrentContentAttributeFacet[model.VideoSource]
}

func (f videoSourceFacet) Values(query.FacetContext) (map[string]string, error) {
	vsrcs := model.VideoSourceValues()
	values := make(map[string]string, len(vsrcs))
	for _, vr := range vsrcs {
		values[vr.String()] = vr.Label()
	}
	return values, nil
}
