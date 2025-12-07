package registry

import (
	"errors"
	"fmt"

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
	bundles ref.Map[bundle.Bundle]
}

func New(bundles ...bundle.Bundle) (*Builder, error) {
	builder := &Builder{
		bundles: ref.NewMap[bundle.Bundle](),
	}

	for _, bundle := range bundles {
		if builder.bundles.Has(bundle.Ref()) {
			return nil, errors.New("plugin already registered")
		}
		builder.bundles.Set(bundle.Ref(), bundle)
	}

	return builder, nil
}

func (r *Builder) AllRefs() []ref.Ref {
	return r.bundles.Refs()
}

func (r *Builder) Resolve(env env.Env, options ...Option) (*Registry, error) {
	plugins := ref.NewMap[plugin.Plugin]()

	for _, bundle := range r.bundles.Values() {
		bundlePlugins, err := bundle.LoadPlugins(env)
		if err != nil {
			return nil, err
		}

		for _, p := range bundlePlugins {
			if plugins.Has(p.Ref()) {
				return nil, fmt.Errorf("plugin already registered: %s", p.Ref())
			}

			plugins.Set(p.Ref(), p)
		}
	}

	config, err := resolveConfig(env, plugins)
	if err != nil {
		return nil, err
	}

	resolver := &resolver{
		plugins:       plugins,
		config:        *config,
		errorRegistry: resolveErrors(plugins),
		i18nProvider:  resolveI18n(plugins),
	}

	for _, option := range options {
		option(resolver)
	}

	return resolver.resolve()
}

func resolveConfig(env env.Env, plugins ref.Map[plugin.Plugin]) (*config_resolver.Resolved, error) {
	configLookup, err := lookup.NewFromEnv(env)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}

	configParams := slice.FlatMap(plugins.Refs(), func(ref ref.Ref) []config_registry.Param {
		return plugins.Get(ref).ConfigParams()
	})

	configResolver := config_resolver.New(configLookup, configParams...)
	resolvedConfig, err := configResolver.Resolve()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", Err, err)
	}
	return &resolvedConfig, nil
}

func resolveErrors(plugins ref.Map[plugin.Plugin]) error_registry.Registry {
	return error_registry.New(slice.Map(plugins.Values(), func(plugin plugin.Plugin) error_registry.Option {
		return error_registry.WithEntries(plugin.Errors())
	})...)
}

func resolveI18n(plugins ref.Map[plugin.Plugin]) i18n.Provider {
	return i18n.Providers(
		i18n.NewProvider(
			slice.FlatMap(plugins.Values(), func(plugin plugin.Plugin) []*i18n.Message {
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
