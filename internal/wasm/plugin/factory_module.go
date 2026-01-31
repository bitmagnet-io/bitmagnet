package plugin

import (
	context "context"
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	plugin_api "github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type module struct {
	api.Module
	getIndexer       func() (plugin_api.Indexer, error)
	getHTTPHandler   func() (plugin_api.HTTPHandler, error)
	getSearchAdapter func() (plugin_api.SearchAdapter, error)
	getTorrentTarget func() (plugin_api.TorrentTarget, error)
}

type instanceBuilder struct {
	env          env.Env
	runtime      wazero.Runtime
	compiled     wazero.CompiledModule
	moduleConfig wazero.ModuleConfig
}

func (i instanceBuilder) newModule(
	ctx context.Context,
) (*module, error) {
	apiModule, err := i.runtime.InstantiateModule(
		ctx,
		i.compiled,
		i.moduleConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	getIndexer := sync.OnceValues(func() (plugin_api.Indexer, error) {
		plugin := plugin_api.IndexerPlugin{}
		return plugin.LoadModule(ctx, apiModule)
	})

	getHTTPHandler := sync.OnceValues(func() (plugin_api.HTTPHandler, error) {
		plugin := plugin_api.HTTPHandlerPlugin{}
		return plugin.LoadModule(ctx, apiModule)
	})

	getSearchAdapter := sync.OnceValues(func() (plugin_api.SearchAdapter, error) {
		plugin := plugin_api.SearchAdapterPlugin{}
		return plugin.LoadModule(ctx, apiModule)
	})

	getTorrentTarget := sync.OnceValues(func() (plugin_api.TorrentTarget, error) {
		plugin := plugin_api.TorrentTargetPlugin{}
		return plugin.LoadModule(ctx, apiModule)
	})

	return &module{
		Module:           apiModule,
		getIndexer:       getIndexer,
		getHTTPHandler:   getHTTPHandler,
		getSearchAdapter: getSearchAdapter,
		getTorrentTarget: getTorrentTarget,
	}, nil
}
