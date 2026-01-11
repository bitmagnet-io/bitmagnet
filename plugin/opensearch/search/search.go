//go:build wasip1

package search

import (
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/search"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func New(client *opensearchapi.Client, indexPrefix string) api.SearchAdapter {
	return &adapter{
		client:      client,
		indexPrefix: indexPrefix,
	}
}

type adapter struct {
	client      *opensearchapi.Client
	indexPrefix string
}

type customAgg struct {
	name string
	body map[string]any
}

func (agg customAgg) Name() string {
	return agg.name
}

func (agg customAgg) Map() map[string]interface{} {
	return agg.body
}

// var facetFields = map[string]string{
// 	"content_type": "contentType",
// }

type responseBucket struct {
	Key      string `json:"key"`
	DocCount int32  `json:"doc_count"`
}

type resultAgg interface {
	resultItems() ([]*search.FacetResultItem, error)
}

type filteredAgg[T any] struct {
	Agg T `json:"agg"`
}

type responseTermAgg filteredAgg[struct {
	Buckets []responseBucket `json:"buckets"`
}]

func (agg responseTermAgg) resultItems() ([]*search.FacetResultItem, error) {
	items := make([]*search.FacetResultItem, 0, len(agg.Agg.Buckets))
	for _, bucket := range agg.Agg.Buckets {
		items = append(items, &search.FacetResultItem{
			Value: bucket.Key,
			Count: bucket.DocCount,
		})
	}
	return items, nil
}

type nestedBucket struct {
	Key        string `json:"key"`
	DocCount   int32  `json:"doc_count"`
	BackToRoot struct {
		DocCount int32 `json:"doc_count"`
	} `json:"back_to_root"`
}

type nestedAgg filteredAgg[struct {
	Value struct {
		Buckets []nestedBucket `json:"buckets"`
	} `json:"value"`
}]

func (agg nestedAgg) resultItems() ([]*search.FacetResultItem, error) {
	items := make([]*search.FacetResultItem, 0, len(agg.Agg.Value.Buckets))
	for _, bucket := range agg.Agg.Value.Buckets {
		items = append(items, &search.FacetResultItem{
			Value: bucket.Key,
			Count: bucket.BackToRoot.DocCount,
		})
	}
	return items, nil
}
