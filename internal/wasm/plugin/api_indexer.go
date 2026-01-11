package plugin

import (
	context "context"

	plugin_api "github.com/bitmagnet-io/bitmagnet/proto/api"
)

type apiIndexer struct {
	*apiService[plugin_api.Indexer]
}

func (api *apiIndexer) Index(ctx context.Context, indexPayload *plugin_api.IndexPayload) (*plugin_api.Empty, error) {
	return &plugin_api.Empty{},
		api.do(ctx, func(indexer plugin_api.Indexer) error {
			_, err := indexer.Index(ctx, indexPayload)
			return err
		})
}
