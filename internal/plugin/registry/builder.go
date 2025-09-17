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
	"github.com/bitmagnet-io/bitmagnet/internal/error_registry"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
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
		Builder:       b,
		config:        *config,
		errorRegistry: b.resolveErrors(),
		i18nProvider:  b.resolveI18n(),
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
		return b.plugins.Get(ref).ConfigParams()
	})

	configResolver := config_resolver.New(configLookup, configParams...)
	resolvedConfig, err := configResolver.Resolve()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}
	return &resolvedConfig, nil
}

func (b *Builder) resolveErrors() error_registry.Registry {
	return error_registry.New(slice.Map(b.plugins.Values(), func(plugin plugin.Plugin) error_registry.Option {
		return error_registry.WithEntries(plugin.Errors())
	})...)
}

func (b *Builder) resolveI18n() i18n.Provider {
	return i18n.Providers(
		i18n.NewProvider(
			slice.FlatMap(b.plugins.Values(), func(plugin plugin.Plugin) []*i18n.Message {
				messages := plugin.I18nMessages()

				messages = append(messages, slice.Map(plugin.Errors().Entries(), func(entry ref.Entry[error]) *i18n.Message {
					return i18n.NewMessage(
						entry.Ref.String(),
						fmt.Sprintf("message for error: %s", entry.Ref),
						i18n.WithOther(entry.Value.Error()),
					)
				})...)

				messages = append(messages, slice.Map(plugin.ConfigParams(), func(param config_registry.Param) *i18n.Message {
					return i18n.NewMessage(
						param.Ref.String(),
						fmt.Sprintf("description for config: %s", param.Ref),
						i18n.WithOther(param.Description()),
					)
				})...)

				return messages
			})...,
		),
	)
}
