package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
)

const QueueJobQueueFacetKey = "queue"

func QueueJobQueueFacet(options ...query.FacetOption) query.Facet {
	return queueJobQueueFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(QueueJobQueueFacetKey),
				query.FacetHasLabel("Queue"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
	}
}

type queueJobQueueFacet struct {
	query.FacetConfig
}

var queueNames = []string{"process_torrent", "process_torrent_batch"}

func (f queueJobQueueFacet) Values(query.FacetContext) (map[string]string, error) {
	values := make(map[string]string)
	for _, n := range queueNames {
		values[n] = n
	}
	return values, nil
}

func (f queueJobQueueFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	values := filter.Values()
	if len(values) == 0 {
		return nil
	}
	return []query.Criteria{
		QueueJobQueueCriteria(filter.Values()...),
	}
}
