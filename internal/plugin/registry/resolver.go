package registry

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	config_resolver "github.com/bitmagnet-io/bitmagnet/internal/config/resolver"
	"github.com/bitmagnet-io/bitmagnet/internal/error_registry"
	"github.com/bitmagnet-io/bitmagnet/internal/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type resolver struct {
	*Builder
	config          config_resolver.Resolved
	i18nProvider    i18n.Provider
	errorRegistry   error_registry.Registry
	enabledPlugins  []ref.Ref
	disabledPlugins []ref.Ref
}

func (r *resolver) resolve() (*Registry, error) {
	resolvedNamesMap, err := r.resolveNames()
	if err != nil {
		return nil, fmt.Errorf("%w: %w: %w", Err, ErrResolve, err)
	}

	resolvedNames := slices.Collect(maps.Values(resolvedNamesMap))

	slices.SortFunc(resolvedNames, func(a, b ref.Ref) int {
		return cmp.Compare(a.String(), b.String())
	})

	requiredBy := make(map[string][]ref.Ref)

	pluginInfos := plugin.PluginInfos(slice.Map(r.AllRefs(), func(ref ref.Ref) plugin.PluginInfo {
		_, enabled := resolvedNamesMap[ref.String()]

		dependsOn := r.DependenciesOf(ref)

		for _, dep := range dependsOn {
			requiredBy[dep.String()] = append(requiredBy[dep.String()], ref)
		}

		return plugin.PluginInfo{
			Ref:       ref,
			Enabled:   enabled,
			DependsOn: dependsOn,
		}
	}))

	for i := range pluginInfos {
		pluginInfos[i].RequiredBy = requiredBy[pluginInfos[i].Ref.String()]
	}

	return &Registry{
		pluginInfos: pluginInfos,
		config:      r.config,
		commands: slice.FlatMap(pluginInfos, func(info plugin.PluginInfo) []plugin.Command {
			if !info.Enabled {
				return nil
			}

			return r.plugins.Get(info.Ref).Commands()
		}),
		fxOption: fx.Options(
			fx.Supply(
				r.config,
				r.i18nProvider,
				r.errorRegistry,
			),
			fx.Options(slice.Map(resolvedNames, func(ref ref.Ref) fx.Option {
				return r.plugins.Get(ref).FXOption()
			})...),
			fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
				l := &fxevent.ZapLogger{Logger: log.Named("fx")}
				l.UseLogLevel(zapcore.DebugLevel)

				return fxLogger{l}
			}),
			fx.Supply(
				pluginInfos,
			),
		),
	}, nil
}

func (r *resolver) resolveNames() (map[string]ref.Ref, error) {
	for _, ref := range r.disabledPlugins {
		if !r.plugins.Has(ref) {
			return nil, fmt.Errorf("%w: %s", ErrUnknownPlugin, ref)
		}
	}

	disabledNamesMap := namesToMap(r.disabledPlugins)
	namesMap := namesToMap(r.enabledPlugins)

	resultMap := make(map[string]ref.Ref)

	var addPlugin func(ref.Ref) error
	addPlugin = func(ref ref.Ref) error {
		if !r.plugins.Has(ref) {
			return fmt.Errorf("%w: %s", ErrUnknownPlugin, ref)
		}

		if _, ok := disabledNamesMap[ref.String()]; ok {
			return fmt.Errorf("%w: %s", ErrDisabled, ref)
		}

		if _, ok := resultMap[ref.String()]; ok {
			return nil
		}

		for _, dep := range r.DependenciesOf(ref) {
			if err := addPlugin(dep); err != nil {
				return fmt.Errorf("%w: %s: %w", ErrDependency, ref, err)
			}
		}

		resultMap[ref.String()] = ref

		return nil
	}

	for _, pl := range r.plugins.Values() {
		enabled := false

		if activationRef, ok := pl.ActivationRef().Value(); !ok {
			enabled = true
		} else {
			if activationParam, ok := r.config.Param(activationRef); ok {
				if activation, ok := activationParam.Value().(plugin.Activation); ok {
					if activation == plugin.ActivationEnabled {
						enabled = true
					}
				}
			}
		}

		if enabled {
			if err := addPlugin(pl.Ref()); err != nil {
				return nil, err
			}
		}
	}

	for _, name := range namesMap {
		if err := addPlugin(name); err != nil {
			return nil, err
		}
	}

	return resultMap, nil
}

func namesToMap(names []ref.Ref) map[string]ref.Ref {
	result := make(map[string]ref.Ref, len(names))

	for _, name := range names {
		result[name.String()] = name
	}

	return result
}

type fxLogger struct {
	fxevent.Logger
}

func (l fxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.Started:
		if errors.Is(e.Err, context.Canceled) {
			return
		}
	case *fxevent.RollingBack:
		if errors.Is(e.StartErr, context.Canceled) {
			return
		}
	}

	l.Logger.LogEvent(event)
}
