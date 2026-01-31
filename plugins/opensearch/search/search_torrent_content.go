package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/proto/common/model"
	"github.com/bitmagnet-io/bitmagnet/proto/common/search"
	"github.com/bitmagnet-io/plugin-opensearch/shared"
	"github.com/defensestation/osquery/v2"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func (a *adapter) SearchTorrentContent(
	ctx context.Context,
	params *search.Params,
) (*search.TorrentContentResult, error) {
	req, err := searchTorrentContentReq(params)
	if err != nil {
		return nil, err
	}
	body, err := req.MarshalJSON()
	if err != nil {
		return nil, err
	}
	result, err := a.client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{a.indexPrefix + string(shared.RecordTorrent)},
		Body:    bytes.NewReader(body),
	})
	if err != nil {
		return nil, err
	}
	var items []*model.TorrentContent
	for i, hit := range result.Hits.Hits {
		var rc model.TorrentContent
		err := json.Unmarshal(hit.Source, &rc)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to unmarshal hit %d: %v", shared.Err, i, err)
		}
		items = append(items, &rc)
	}
	totalCount := int32(result.Hits.Total.Value)
	var rawFacets map[string]json.RawMessage
	if len(result.Aggregations) > 0 {
		err = json.Unmarshal(result.Aggregations, &rawFacets)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to unmarshal aggregations: %v", shared.Err, err)
		}
	} else {
		rawFacets = make(map[string]json.RawMessage)
	}
	facetResult, err := allFacets.res(rawFacets)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to process facets: %v", shared.Err, err)
	}
	return &search.TorrentContentResult{
		TotalCount: &totalCount,
		Facets:     facetResult,
		Items:      items,
	}, nil
}

func searchTorrentContentReq(params *search.Params) (*osquery.SearchRequest, error) {
	req := osquery.Search()
	filters, err := createCriteriaFilters(params)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create criteria filters: %v", shared.Err, err)
	}
	if params.QueryString != nil && *params.QueryString != "" {
		filters = append(filters,
			osquery.MultiMatch(*params.QueryString).
				Fields(
					"torrent.name",
					"content.title",
					"content.originalTitle",
				),
		)
	}

	if len(filters) > 0 {
		req.Query(osquery.Bool().Must(filters...))
	}

	if params.Limit != nil {
		size := uint64(*params.Limit)
		req.Size(size)
	}
	offset := uint64(0)
	if params.Offset != nil {
		offset += uint64(*params.Offset)
	}
	if params.Page != nil && params.Limit != nil {
		offset += uint64((*params.Page - 1) * (*params.Limit))
	}
	if offset > 0 {
		req.From(offset)
	}
	allFacets.req(req, params)

	req.SourceExcludes("torrent.files")

	return req, nil
}
