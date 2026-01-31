package multi

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Index struct {
	ref.Ref
	Name    string
	Adapter search.Base
}

type IndexInfo struct {
	Ref         ref.Ref
	Name        string
	ResultTypes []search.ResultType
}

func (a Index) Info() IndexInfo {
	resultTypes := make([]search.ResultType, 0)
	if _, ok := any(a.Adapter).(search.TorrentContent); ok {
		resultTypes = append(resultTypes, search.ResultTypeTorrentContent)
	}

	if _, ok := any(a.Adapter).(search.TorrentFiles); ok {
		resultTypes = append(resultTypes, search.ResultTypeTorrentFile)
	}

	if _, ok := any(a.Adapter).(search.Content); ok {
		resultTypes = append(resultTypes, search.ResultTypeContent)
	}

	return IndexInfo{
		Ref:         a.Ref,
		Name:        a.Name,
		ResultTypes: resultTypes,
	}
}

type Search struct {
	search.Base
	defaultIndex ref.Ref
	adapters     ref.Map[Index]
}

func New(adapters ...Index) Search {
	a := Search{
		adapters: ref.NewMap[Index](),
	}

	for i, adapter := range adapters {
		if i == 0 {
			a.defaultIndex = adapter.Ref
		}

		a.adapters.Set(adapter.Ref, adapter)
	}

	return a
}

func (multi Search) DefaultIndex() ref.Ref {
	return multi.defaultIndex
}

func (multi Search) Indexes() []IndexInfo {
	return slice.Map(multi.adapters.Values(), func(adapter Index) IndexInfo {
		return adapter.Info()
	})
}

func (multi Search) TorrentContent(ctx context.Context, params search.Params) (search.TorrentContentResult, error) {
	adapter, err := getIndex[search.TorrentContent](multi, params.Index)
	if err != nil {
		return search.TorrentContentResult{}, err
	}

	return adapter.TorrentContent(ctx, params)
}

func (multi Search) TorrentFiles(ctx context.Context, params search.Params) (search.TorrentFilesResult, error) {
	adapter, err := getIndex[search.TorrentFiles](multi, params.Index)
	if err != nil {
		return search.TorrentFilesResult{}, err
	}

	return adapter.TorrentFiles(ctx, params)
}

func (multi Search) Content(ctx context.Context, params search.Params) (search.ContentResult, error) {
	adapter, err := getIndex[search.Content](multi, params.Index)
	if err != nil {
		return search.ContentResult{}, err
	}

	return adapter.Content(ctx, params)
}

func getIndex[T any](multi Search, index ref.Nullable) (T, error) {
	resolvedRef := multi.defaultIndex
	if index.Valid {
		resolvedRef = index.Ref
	}

	untyped, ok := multi.adapters.GetOK(resolvedRef)
	if !ok {
		var zero T
		return zero, fmt.Errorf("%w: %w: %s", search.Err, search.ErrUnknownIndex, resolvedRef)
	}

	adapter, ok := any(untyped.Adapter).(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("%w: %w: %s", search.Err, search.ErrInvalidAdapterType, resolvedRef)
	}

	return adapter, nil
}
