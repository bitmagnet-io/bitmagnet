package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"gorm.io/gen"
)

const TorrentSourceFacetKey = "torrent_source"

func TorrentSourceFacet(options ...query.FacetOption) query.Facet {
	return torrentSourceFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(TorrentSourceFacetKey),
				query.FacetHasLabel("Torrent Source"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
	}
}

type torrentSourceFacet struct {
	query.FacetConfig
}

func (torrentSourceFacet) Values(ctx query.FacetContext) (map[string]string, error) {
	q := ctx.Query().TorrentSource

	sources, sourcesErr := q.WithContext(ctx.Context()).Find()
	if sourcesErr != nil {
		return nil, sourcesErr
	}

	values := make(map[string]string, len(sources))
	for _, s := range sources {
		values[s.Key] = s.Name
	}

	return values, nil
}

func (torrentSourceFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	if len(filter) == 0 {
		return []query.Criteria{}
	}

	return []query.Criteria{
		TorrentSourceCriteria(filter.Values()...),
	}
}

func TorrentSourceCriteria(keys ...string) query.Criteria {
	return query.GenCriteria(func(ctx query.DBContext) (query.Criteria, error) {
		q := ctx.Query()

		return query.RawCriteria{
			Query: gen.Exists(
				q.TorrentsTorrentSource.Where(
					q.TorrentsTorrentSource.InfoHash.EqCol(q.TorrentContent.InfoHash),
					q.TorrentsTorrentSource.Source.In(keys...),
				),
			),
		}, nil
	})
}
