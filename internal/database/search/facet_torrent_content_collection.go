package search

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/query"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"strings"
)

type torrentContentCollectionFacet struct {
	query.FacetConfig
	collectionType string
}

func (f torrentContentCollectionFacet) Values(ctx query.FacetContext) (map[string]string, error) {
	q := ctx.Query().ContentCollection
	colls, err := ctx.Query().ContentCollection.WithContext(ctx.Context()).Where(
		q.Type.Eq(f.collectionType),
	).Find()
	if err != nil {
		return nil, err
	}
	values := make(map[string]string, len(colls))
	for _, coll := range colls {
		values[coll.Source+":"+coll.ID] = coll.Name
	}
	return values, nil
}

func (f torrentContentCollectionFacet) Criteria(filter query.FacetFilter) []query.Criteria {
	sourceMap := make(map[string]map[string]struct{})
	for _, value := range filter.Values() {
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			continue
		}
		source, id := parts[0], parts[1]
		if _, ok := sourceMap[source]; !ok {
			sourceMap[source] = make(map[string]struct{})
		}
		sourceMap[source][id] = struct{}{}
	}
	criteria := make([]query.Criteria, 0, len(sourceMap))
	for source, idMap := range sourceMap {
		refs := make([]model.ContentCollectionRef, 0, len(idMap))
		for id := range idMap {
			refs = append(refs, model.ContentCollectionRef{
				Type:   f.collectionType,
				Source: source,
				ID:     id,
			})
		}
		switch f.Logic() {
		case model.FacetLogicOr:
			criteria = append(criteria, ContentCollectionCriteria(refs...))
		case model.FacetLogicAnd:
			for _, ref := range refs {
				criteria = append(criteria, ContentCollectionCriteria(ref))
			}
		}
	}
	return criteria
}
