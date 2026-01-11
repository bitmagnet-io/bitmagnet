package plugin

import (
	context "context"

	plugin_api "github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/search"
)

type apiSearchAdapter struct {
	*apiService[plugin_api.SearchAdapter]
}

func (api *apiSearchAdapter) SearchTorrentContent(ctx context.Context, params *search.Params) (*search.TorrentContentResult, error) {
	var result *search.TorrentContentResult

	err := api.do(ctx, func(searchAdapter plugin_api.SearchAdapter) error {
		var err error
		result, err = searchAdapter.SearchTorrentContent(ctx, params)
		return err
	})

	return result, err
}

func (api *apiSearchAdapter) SearchTorrentFiles(ctx context.Context, params *search.Params) (*search.TorrentFilesResult, error) {
	var result *search.TorrentFilesResult
	err := api.do(ctx, func(searchAdapter plugin_api.SearchAdapter) error {
		var err error
		result, err = searchAdapter.SearchTorrentFiles(ctx, params)
		return err
	})
	return result, err
}
