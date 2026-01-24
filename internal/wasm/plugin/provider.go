package plugin

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/param"
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/spf13/afero"
	"github.com/tetratelabs/wazero"
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

func (w wasmProvider) loadPlugin(env env.Env, plugin pluginPathAlias) (*Plugin, error) {
	// todo: FS
	dir := afero.NewBasePathFs(env.FSRoot(), plugin.path)

	manifestBytes, err := afero.ReadFile(dir, "manifest.json")
	if err != nil {
		return nil, err
	}

	manifest, err := ParseManifest(manifestBytes)
	if err != nil {
		return nil, err
	}

	name := plugin.alias
	if name == "" {
		name = manifest.Name
	}

	ref, err := ref.Parse(name)
	if err != nil {
		return nil, err
	}

	configParams := make([]config_registry.Param, 0, len(manifest.Config))
	for key, schema := range manifest.Config {
		paramRef, err := ref.Sub(key)
		if err != nil {
			return nil, err
		}

		param, err := param.New(param.JSONSchemaDecoder[any](schema))
		if err != nil {
			return nil, err
		}

		configParams = append(configParams, config_registry.Param{
			Ref:     paramRef,
			Plugin:  ref,
			Untyped: param,
		})
	}

	data, err := afero.ReadFile(dir, "plugin.wasm")
	if err != nil {
		return nil, err
	}

	return &Plugin{
		ref:              ref,
		manifest:         manifest,
		configParams:     configParams,
		data:             data,
		compilationCache: w.compilationCache,
	}, nil
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
