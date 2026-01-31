package search

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/bitmagnet-io/bitmagnet/proto/common/model"
	"github.com/bitmagnet-io/bitmagnet/proto/common/search"
	"github.com/bitmagnet-io/plugin-opensearch/shared"
	"github.com/defensestation/osquery/v2"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func (a *adapter) SearchTorrentFiles(
	ctx context.Context,
	params *search.Params,
) (*search.TorrentFilesResult, error) {
	req, err := searchTorrentFilesReq(params)
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
	var aggs searchTorrentFilesResultAggs
	err = json.Unmarshal(result.Aggregations, &aggs)
	if err != nil {
		return nil, err
	}
	var items []*model.TorrentFile
	for _, hit := range aggs.Files.Files.Hits.Hits {
		items = append(items, hit.Source)
	}
	var rawFacets map[string]json.RawMessage
	err = json.Unmarshal(result.Aggregations, &rawFacets)
	if err != nil {
		return nil, err
	}
	facetResult, err := allFacets.res(rawFacets)
	if err != nil {
		return nil, err
	}
	totalCount := int32(aggs.Files.Files.Hits.Total.Value)
	return &search.TorrentFilesResult{
		TotalCount: &totalCount,
		Facets:     facetResult,
		Items:      items,
	}, nil
}

type searchTorrentFilesResultAggs struct {
	Files struct {
		Files struct {
			Hits struct {
				Total struct {
					Value int `json:"value"`
				} `json:"total"`
				Hits []struct {
					Source *model.TorrentFile `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		} `json:"files"`
	} `json:"files"`
}

func searchTorrentFilesReq(params *search.Params) (*osquery.SearchRequest, error) {
	req := osquery.Search()
	filters, err := createCriteriaFilters(params)
	if err != nil {
		return nil, err
	}

	if len(filters) > 0 {
		req.Query(osquery.Bool().Must(filters...))
	}

	req.Size(1)

	allFacets.req(req, params)

	topHits := &topHitsCustom{
		name:      "files",
		sort:      "torrent.files.index",
		direction: "asc",
		missing:   0,
	}
	if params.Limit != nil {
		topHits.size = uint64(*params.Limit)
	}
	offset := uint64(0)
	if params.Offset != nil {
		offset += uint64(*params.Offset)
	}
	if params.Page != nil && params.Limit != nil {
		offset += uint64((*params.Page - 1) * (*params.Limit))
	}
	if offset > 0 {
		topHits.from = offset
	}

	req.Aggs(
		osquery.NestedAgg("files", "torrent.files").Aggs(topHits),
	)

	req.SourceExcludes("torrent.files")

	return req, nil
}

type topHitsCustom struct {
	name      string
	sort      string
	direction string
	missing   any
	size      uint64
	from      uint64
}

func (m *topHitsCustom) Name() string {
	return m.name
}

func (m *topHitsCustom) Map() map[string]any {
	return map[string]any{
		"top_hits": map[string]any{
			"size": m.size,
			"from": m.from,
			"sort": []map[string]any{
				{
					m.sort: map[string]any{
						"order":   m.direction,
						"missing": m.missing,
					},
				},
			},
		},
	}
}
