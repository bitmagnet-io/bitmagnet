package search

import (
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const TorrentFileTypeFacetKey = "file_type"

func TorrentFileTypeFacet(options ...query.FacetOption) query.Facet {
	return torrentFileTypeFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(TorrentFileTypeFacetKey),
				query.FacetHasLabel("File Type"),
				query.FacetUsesOrLogic(),
				query.FacetHasAggregationOption(query.RequireJoin(model.TableNameTorrentContent)),
				query.FacetTriggersCte(),
			}, options...)...,
		),
	}
}

type torrentFileTypeFacet struct {
	query.FacetConfig
}

func (torrentFileTypeFacet) Values(query.FacetContext) (map[string]string, error) {
	fts := model.FileTypeValues()
	values := make(map[string]string, len(fts))
	for _, vr := range fts {
		values[vr.String()] = vr.Label()
	}
	return values, nil
}

func (f torrentFileTypeFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	return []query.Criteria{query.GenCriteria(func(query.DBContext) (query.Criteria, error) {
		if len(filter) == 0 {
			return query.AndCriteria{}, nil
		}
		fileTypes := make([]model.FileType, 0, len(filter))
		for _, v := range filter.Values() {
			ft, ftErr := model.ParseFileType(v)
			if ftErr != nil {
				return nil, errors.New("invalid file type filter specified")
			}
			fileTypes = append(fileTypes, ft)
		}
		return TorrentFileTypeCriteria(fileTypes...), nil
	})}
}
