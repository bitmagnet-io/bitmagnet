package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const TorrentTagFacetKey = "torrent_tag"

func TorrentTagsFacet(options ...query.FacetOption) query.Facet {
	return torrentTagFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(TorrentTagFacetKey),
				query.FacetHasLabel("Torrent Tag"),
				query.FacetUsesAndLogic(),
				query.FacetHasAggregationOption(query.RequireJoin(model.TableNameTorrentContent)),
			}, options...)...,
		),
	}
}

type torrentTagFacet struct {
	query.FacetConfig
}

func (torrentTagFacet) Values(ctx query.FacetContext) (map[string]string, error) {
	q := ctx.Query().TorrentTag
	tags, tagsErr := q.WithContext(ctx.Context()).Distinct(q.Name).Find()
	if tagsErr != nil {
		return nil, tagsErr
	}
	values := make(map[string]string, len(tags))
	for _, tag := range tags {
		values[tag.Name] = tag.Name
	}
	return values, nil
}

func (f torrentTagFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	criteria := make([]query.Criteria, len(filter))
	for i, tag := range filter.Values() {
		criteria[i] = TorrentTagCriteria(tag)
	}
	return criteria
}
