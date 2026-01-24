package plugin

import (
	context "context"
	"fmt"
	"sync"

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

func (i instanceBuilder) newModule(
	ctx context.Context,
) (*module, error) {
	runtime := wazero.NewRuntimeWithConfig(ctx, i.runtimeConfig)

	for _, instantiator := range i.instantiators {
		if err := instantiator(ctx, runtime); err != nil {
			return nil, fmt.Errorf("failed to instantiate host module: %w", err)
		}
	}

	compiled, err := runtime.CompileModule(ctx, i.data)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	apiModule, err := runtime.InstantiateModule(
		ctx,
		compiled,
		i.moduleConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	var (
		getIndexer       func() (plugin_api.Indexer, error)
		getHTTPHandler   func() (plugin_api.HTTPHandler, error)
		getSearchAdapter func() (plugin_api.SearchAdapter, error)
		getTorrentTarget func() (plugin_api.TorrentTarget, error)
	)

	if i.manifest.Capabilities.Indexer != nil {
		getIndexer = sync.OnceValues(func() (plugin_api.Indexer, error) {
			plugin := plugin_api.IndexerPlugin{}
			return plugin.LoadModule(ctx, apiModule)
		})
	}

	if i.manifest.Capabilities.HTTPHandler != nil {
		getHTTPHandler = sync.OnceValues(func() (plugin_api.HTTPHandler, error) {
			plugin := plugin_api.HTTPHandlerPlugin{}
			return plugin.LoadModule(ctx, apiModule)
		})
	}

	if i.manifest.Capabilities.SearchAdapter != nil {
		getSearchAdapter = sync.OnceValues(func() (plugin_api.SearchAdapter, error) {
			plugin := plugin_api.SearchAdapterPlugin{}
			return plugin.LoadModule(ctx, apiModule)
		})
	}

	if i.manifest.Capabilities.TorrentTarget != nil {
		getTorrentTarget = sync.OnceValues(func() (plugin_api.TorrentTarget, error) {
			plugin := plugin_api.TorrentTargetPlugin{}
			return plugin.LoadModule(ctx, apiModule)
		})
	}

	return &module{
		Module:           apiModule,
		getIndexer:       getIndexer,
		getHTTPHandler:   getHTTPHandler,
		getSearchAdapter: getSearchAdapter,
		getTorrentTarget: getTorrentTarget,
	}, nil
}
