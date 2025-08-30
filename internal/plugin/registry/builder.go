package registry

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/config/lookup"
	config_registry "github.com/bitmagnet-io/bitmagnet/internal/config/registry"
	config_resolver "github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/env"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/bundle"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type Builder struct {
	plugins ref.Map[plugin.Plugin]
}

func New(bundles ...bundle.Bundle) (*Builder, error) {
	builder := &Builder{
		plugins: ref.NewMap[plugin.Plugin](),
	}

	for _, bundle := range bundles {
		for _, plugin := range bundle.Plugins() {
			if builder.plugins.Has(plugin.Ref()) {
				return nil, errors.New("plugin already registered")
			}
			builder.plugins.Set(plugin.Ref(), plugin)
		}
	}

	return builder, nil
}

func (b *Builder) AllRefs() []ref.Ref {
	return b.plugins.Refs()
}

func (b *Builder) DependenciesOf(rf ref.Ref) []ref.Ref {
	plugin, ok := b.plugins.GetOK(rf)
	if !ok {
		return nil
	}

	var result []ref.Ref

	seen := make(map[string]struct{})

	for _, dependency := range plugin.Dependencies() {
		if _, ok := seen[dependency.String()]; !ok {
			seen[dependency.String()] = struct{}{}

			result = append(result, dependency)
			for _, subDependency := range b.DependenciesOf(dependency) {
				if _, ok = seen[subDependency.String()]; !ok {
					seen[subDependency.String()] = struct{}{}

					result = append(result, subDependency)
				}
			}
		}
		// todo: Detect circular
	}

	slices.SortFunc(result, func(a, b ref.Ref) int {
		return cmp.Compare(a.String(), b.String())
	})

	return result
}

func (b *Builder) Resolve(env env.Env, options ...Option) (*Registry, error) {
	config, err := b.resolveConfig(env)
	if err != nil {
		return nil, err
	}

	resolver := &resolver{
		Builder:  b,
		Resolved: *config,
	}

	for _, option := range options {
		option(resolver)
	}

	return resolver.resolve()
}

func (b *Builder) resolveConfig(env env.Env) (*config_resolver.Resolved, error) {
	configLookup, err := lookup.NewFromEnv(env)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}

	allRefs := b.AllRefs()

	configParams := slice.FlatMap(allRefs, func(ref ref.Ref) []config_registry.Param {
		return b.plugins.Get(ref).Params()
	})

	configResolver := config_resolver.New(configLookup, configParams...)
	resolvedConfig, err := configResolver.Resolve()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}
	return &resolvedConfig, nil
}
