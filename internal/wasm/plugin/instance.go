package plugin

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/plugin"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type Instance interface {
	Contract() *plugin.Contract
	Close(ctx context.Context)
	Indexer() api.Indexer
	HTTPHandler() api.HTTPHandler
	SearchAdapter() api.SearchAdapter
	TorrentTarget() api.TorrentTarget
}

type instance struct {
	ref      ref.Ref
	contract *plugin.Contract
	pool     *pool.ObjectPool
}

func (i *instance) Ref() ref.Ref {
	return i.ref
}

func (i *instance) Contract() *plugin.Contract {
	return i.contract
}

func (i *instance) Close(ctx context.Context) {
	i.pool.Close(ctx)
}

func (i *instance) Indexer() api.Indexer {
	return &apiIndexer{
		apiService: &apiService[api.Indexer]{
			pool: i.pool,
			getService: func(m *module) (api.Indexer, error) {
				return m.getIndexer()
			},
		},
	}
}

func (i *instance) HTTPHandler() api.HTTPHandler {
	return &apiHTTPHandler{
		apiService: &apiService[api.HTTPHandler]{
			pool: i.pool,
			getService: func(m *module) (api.HTTPHandler, error) {
				return m.getHTTPHandler()
			},
		},
	}
}

func (i *instance) SearchAdapter() api.SearchAdapter {
	return &apiSearchAdapter{
		apiService: &apiService[api.SearchAdapter]{
			pool: i.pool,
			getService: func(m *module) (api.SearchAdapter, error) {
				return m.getSearchAdapter()
			},
		},
	}
}

func (i *instance) TorrentTarget() api.TorrentTarget {
	return &apiTorrentTarget{
		apiService: &apiService[api.TorrentTarget]{
			pool: i.pool,
			getService: func(m *module) (api.TorrentTarget, error) {
				return m.getTorrentTarget()
			},
		},
	}
}
