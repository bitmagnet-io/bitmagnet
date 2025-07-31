package registry

import (
	"cmp"
	"errors"
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/bundle"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Builder struct {
	plugins map[string]plugin.Plugin
}

func New(bundles ...bundle.Bundle) (*Builder, error) {
	builder := &Builder{
		plugins: make(map[string]plugin.Plugin),
	}

	for _, bundle := range bundles {
		for _, plugin := range bundle.Plugins() {
			strPluginRef := plugin.Ref().String()
			if _, ok := builder.plugins[strPluginRef]; ok {
				return nil, errors.New("plugin already registered")
			}
			builder.plugins[strPluginRef] = plugin
		}
	}

	return builder, nil
}

func (r *Builder) AllRefs() []ref.Ref {
	result := slice.Map(slices.Collect(maps.Values(r.plugins)), func(p plugin.Plugin) ref.Ref {
		return p.Ref()
	})

	slices.SortFunc(result, func(a, b ref.Ref) int {
		return cmp.Compare(a.String(), b.String())
	})

	return result
}

func (r *Builder) DependenciesOf(name ref.Ref) []ref.Ref {
	plugin, ok := r.plugins[name.String()]
	if !ok {
		return nil
	}

	var result []ref.Ref
	for _, dependency := range plugin.Dependencies() {
		result = append(result, dependency)
		result = append(result, r.DependenciesOf(dependency)...)
		// todo: Detect circular
	}
	return result
}

func (r *Builder) Resolve(options ...Option) (*Registry, error) {
	resolver := &resolver{
		Builder: r,
	}

	for _, option := range options {
		option(resolver)
	}

	return resolver.resolve()
}
