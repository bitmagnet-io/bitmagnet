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
	return video3dFacet{torrentContentAttributeFacet[model.Video3d]{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(Video3dFacetKey),
				query.FacetHasLabel("Video 3D"),
				query.FacetUsesOrLogic(),
				query.FacetTriggersCte(),
			}, options...)...,
		),
		field: video3dField,
		parse: model.ParseVideo3d,
	}}
}

type video3dFacet struct {
	torrentContentAttributeFacet[model.Video3d]
}

func (f video3dFacet) Values(query.FacetContext) (map[string]string, error) {
	v3ds := model.Video3dValues()
	values := make(map[string]string, len(v3ds))
	for _, vr := range v3ds {
		values[vr.String()] = vr.Label()
	}
	return values, nil
}
