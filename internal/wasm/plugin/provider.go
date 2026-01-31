package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/wasm/host/http_client"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/json_schema"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/bitmagnet-io/bitmagnet/proto/api"
	"github.com/bitmagnet-io/bitmagnet/proto/common/http"
	proto_plugin "github.com/bitmagnet-io/bitmagnet/proto/common/plugin"
	"github.com/spf13/afero"
	"github.com/tetratelabs/wazero"
	wasi_snapshot_preview1 "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type wasmProvider struct {
	compilationCache wazero.CompilationCache
	plugins          []pluginPathAlias
}

func (w wasmProvider) LoadPlugins(env env.Env) ([]plugin.Plugin, error) {
	var plugins []plugin.Plugin
	for _, p := range w.plugins {
		loadedPlugin, err := w.loadPlugin(env, p)
		if err != nil {
			return nil, err
		}

		plugins = append(plugins, loadedPlugin)
	}

	return plugins, nil
}

var pluginPlugin = api.PluginPlugin{}

func (w wasmProvider) loadPlugin(env env.Env, plugin pluginPathAlias) (*Plugin, error) {
	data, err := afero.ReadFile(env.FSRoot(), plugin.path)
	if err != nil {
		return nil, err
	}

	runtime := wazero.NewRuntimeWithConfig(env, wazero.NewRuntimeConfig().
		WithCloseOnContextDone(true).
		WithCompilationCache(w.compilationCache))

	httpEgress := atomic.NewValue[[]*http.Egress](nil)

	for _, instantiator := range []func(ctx context.Context, runtime wazero.Runtime) error{
		func(ctx context.Context, runtime wazero.Runtime) error {
			_, err := wasi_snapshot_preview1.Instantiate(ctx, runtime)
			return err
		},
		http_client.Instantiator(httpEgress),
	} {
		if err := instantiator(env, runtime); err != nil {
			return nil, fmt.Errorf("failed to instantiate host module: %w", err)
		}
	}

	compiled, err := runtime.CompileModule(env, data)
	if err != nil {
		return nil, fmt.Errorf("failed to compile module: %w", err)
	}

	builder := &instanceBuilder{
		env:      env,
		runtime:  runtime,
		compiled: compiled,
		moduleConfig: wazero.NewModuleConfig().
			WithStartFunctions("_initialize"),
	}

	apiModule, err := builder.newModule(env)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}

	pluginAPI, err := pluginPlugin.LoadModule(env, apiModule)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin lifecycle: %w", err)
	}

	ctx, cancel := context.WithTimeout(env, time.Second)
	defer cancel()

	identity, err := pluginAPI.Identify(ctx, &api.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to identify plugin: %w", err)
	}

	name := plugin.alias
	if name == "" {
		name = identity.GetName()
	}

	ref, err := ref.Parse(name)
	if err != nil {
		return nil, err
	}

	defaultContent, err := pluginAPI.Localize(ctx, &api.LocalizeParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin content: %w", err)
	}

	configParams, err := transformConfigParams(
		ref,
		identity.GetConfigParams(),
		defaultContent.GetConfigParams(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to transform config params: %w", err)
	}

	builder.moduleConfig = builder.moduleConfig.
		WithStdout(env).
		WithStderr(env.Stderr()).
		WithSysWalltime().
		WithSysNanotime().
		WithSysNanosleep()

	return &Plugin{
		ref:          ref,
		api:          pluginAPI,
		newModule:    builder.newModule,
		configParams: configParams,
		httpEgress:   httpEgress,
	}, nil
}

func transformConfigParams(
	ref ref.Ref,
	protoParams []*proto_plugin.ConfigParam,
	protoParamsLabels []*proto_plugin.ConfigParamLocalizedContent,
) ([]config_registry.Param, error) {
	configParams := make([]config_registry.Param, 0, len(protoParams))

	for _, protoParam := range protoParams {
		paramRef, err := ref.Sub(protoParam.GetName())
		if err != nil {
			return nil, err
		}

		var schema json_schema.JSONSchema
		if err := json.Unmarshal(protoParam.GetSchema(), &schema); err != nil {
			return nil, fmt.Errorf(
				"failed to unmarshal config param schema for %q: %w",
				protoParam.GetName(),
				err,
			)
		}

		options := []param.Option[any]{param.JSONSchemaDecoder[any](schema)}

		if paramContent, ok := slice.Find(
			protoParamsLabels,
			func(protoLabel *proto_plugin.ConfigParamLocalizedContent) bool {
				return protoLabel.GetName() == protoParam.GetName()
			},
		); ok {
			options = append(options, param.Description[any](paramContent.GetDescription()))
		}

		param, err := param.New(options...)
		if err != nil {
			return nil, err
		}

		configParams = append(configParams, config_registry.Param{
			Ref:     paramRef,
			Plugin:  ref,
			Untyped: param,
		})
	}

	return configParams, nil
}

type pluginPathAlias struct {
	path  string
	alias string
}

type ProviderOption func(*wasmProvider) error

func NewProvider(options ...ProviderOption) (plugin.Provider, error) {
	b := wasmProvider{
		compilationCache: wazero.NewCompilationCache(),
	}

	for _, option := range options {
		if err := option(&b); err != nil {
			return wasmProvider{}, err
		}
	}

	return b, nil
}

func LoadPlugin(path, alias string) ProviderOption {
	return func(b *wasmProvider) error {
		b.plugins = append(b.plugins, pluginPathAlias{
			path:  path,
			alias: alias,
		})

		return nil
	}
}
