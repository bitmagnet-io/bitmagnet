package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"gorm.io/gen/field"
)

const QueueJobStatusFacetKey = "status"

func QueueJobStatusFacet(options ...query.FacetOption) query.Facet {
	return queueJobStatusFacet{
		FacetConfig: query.NewFacetConfig(
			append([]query.FacetOption{
				query.FacetHasKey(QueueJobStatusFacetKey),
				query.FacetHasLabel("Status"),
				query.FacetUsesOrLogic(),
			}, options...)...,
		),
	}
}

type queueJobStatusFacet struct {
	query.FacetConfig
}

func (f queueJobStatusFacet) Values(query.FacetContext) (map[string]string, error) {
	values := make(map[string]string)
	for _, n := range model.QueueJobStatusNames() {
		values[n] = n
	}
	return values, nil
}

func (f queueJobStatusFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	values := filter.Values()
	if len(values) == 0 {
		return nil
	}
	return []query.Criteria{
		query.DaoCriteria{
			Conditions: func(ctx query.DbContext) ([]field.Expr, error) {
				q := ctx.Query()
				return []field.Expr{
					q.QueueJob.Status.In(values...),
				}, nil
			},
		},
	}
}
