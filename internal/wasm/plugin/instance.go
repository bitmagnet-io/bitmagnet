package plugin

import (
	"context"

	plugin_api "github.com/bitmagnet-io/bitmagnet/proto/api"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type Instance interface {
	Close(ctx context.Context)
	Indexer() plugin_api.Indexer
	HTTPHandler() plugin_api.HTTPHandler
	SearchAdapter() plugin_api.SearchAdapter
}

type instance struct {
	manifest Manifest
	pool     *pool.ObjectPool
}

func (i *instance) Close(ctx context.Context) {
	i.pool.Close(ctx)
}

func (i *instance) Indexer() plugin_api.Indexer {
	if i.manifest.Capabilities.Indexer == nil {
		return nil
	}

	return &apiIndexer{
		apiService: &apiService[plugin_api.Indexer]{
			pool: i.pool,
			getService: func(m *module) (plugin_api.Indexer, error) {
				return m.getIndexer()
			},
		},
	}
}

func (i *instance) HTTPHandler() plugin_api.HTTPHandler {
	if i.manifest.Capabilities.HTTPHandler == nil {
		return nil
	}

	return &apiHTTPHandler{
		apiService: &apiService[plugin_api.HTTPHandler]{
			pool: i.pool,
			getService: func(m *module) (plugin_api.HTTPHandler, error) {
				return m.getHTTPHandler()
			},
		},
	}
}

func (i *instance) SearchAdapter() plugin_api.SearchAdapter {
	if i.manifest.Capabilities.SearchAdapter == nil {
		return nil
	}

	return &apiSearchAdapter{
		apiService: &apiService[plugin_api.SearchAdapter]{
			pool: i.pool,
			getService: func(m *module) (plugin_api.SearchAdapter, error) {
				return m.getSearchAdapter()
			},
		},
	}
}
