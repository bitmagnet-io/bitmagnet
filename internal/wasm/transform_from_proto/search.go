package transform_from_proto

import (
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	proto_model "github.com/bitmagnet-io/bitmagnet/proto/common/model"
	proto_search "github.com/bitmagnet-io/bitmagnet/proto/common/search"
)

func SearchTorrentContentResult(
	result *proto_search.TorrentContentResult,
) search.TorrentContentResult {
	var (
		totalCount           model.NullUint
		totalCountIsEstimate bool
		hasNextPage          model.NullBool
	)
	if t := result.GetTotalCount(); t != nil {
		totalCount = model.NullUint{
			Uint:  uint(*t),
			Valid: true,
		}
	}
	if e := result.GetTotalCountIsEstimate(); e != nil {
		totalCountIsEstimate = *e
	}
	if h := result.GetHasNextPage(); h != nil {
		hasNextPage = model.NullBool{
			Bool:  *h,
			Valid: true,
		}
	}
	return search.TorrentContentResult{
		TotalCount:           totalCount,
		TotalCountIsEstimate: totalCountIsEstimate,
		HasNextPage:          hasNextPage,
		Items: slice.Map(result.GetItems(), func(tc *proto_model.TorrentContent) search.TorrentContentResultItem {
			return search.TorrentContentResultItem{
				TorrentContent: TorrentContent(tc),
			}
		}),
		Facets: slice.Map(result.GetFacets(), SearchFacetResult),
	}
}

func SearchTorrentFilesResult(
	result *proto_search.TorrentFilesResult,
) search.TorrentFilesResult {
	var (
		totalCount           model.NullUint
		totalCountIsEstimate bool
		hasNextPage          model.NullBool
	)
	if t := result.GetTotalCount(); t != nil {
		totalCount = model.NullUint{
			Uint:  uint(*t),
			Valid: true,
		}
	}
	if e := result.GetTotalCountIsEstimate(); e != nil {
		totalCountIsEstimate = *e
	}
	if h := result.GetHasNextPage(); h != nil {
		hasNextPage = model.NullBool{
			Bool:  *h,
			Valid: true,
		}
	}
	return search.TorrentFilesResult{
		TotalCount:           totalCount,
		TotalCountIsEstimate: totalCountIsEstimate,
		HasNextPage:          hasNextPage,
		Items: slice.Map(result.GetItems(), func(tc *proto_model.TorrentFile) model.TorrentFile {
			return TorrentFile(tc)
		}),
		Facets: slice.Map(result.GetFacets(), SearchFacetResult),
	}
}

func SearchFacetResult(
	result *proto_search.FacetResult,
) search.FacetResult {
	return search.FacetResult{
		Key:   search.Facet(result.GetKey()),
		Logic: model.FacetLogic(result.GetLogic()),
		Items: slice.Map(result.GetItems(), func(item *proto_search.FacetResultItem) search.FacetResultItem {
			return search.FacetResultItem{
				Value: item.GetValue(),
				Label: func() string {
					label := item.GetLabel()
					if label != nil {
						return *label
					}
					return item.GetValue()
				}(),
				Count: uint(item.GetCount()),
				IsEstimate: func() bool {
					value := item.GetIsEstimate()
					return value != nil && *value
				}(),
			}
		}),
	}
}
