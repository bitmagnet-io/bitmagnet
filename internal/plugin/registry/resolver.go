package registry

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

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
	defaultPlugins  bool
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

	enabledPlugins := slice.Map(resolvedNames, func(name ref.Ref) plugin.Plugin {
		return r.plugins[name.String()]
	})

	pluginInfos := PluginInfos(slice.Map(r.AllRefs(), func(ref ref.Ref) PluginInfo {
		_, enabled := resolvedNamesMap[ref.String()]

		return PluginInfo{
			Ref:       ref,
			Enabled:   enabled,
			DependsOn: r.DependenciesOf(ref),
		}
	}))

	return &Registry{
		pluginInfos: pluginInfos,
		fxOption: fx.Options(
			fx.Options(slice.Map(enabledPlugins, func(plugin plugin.Plugin) fx.Option {
				return plugin.FXOption()
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
		commands: slice.FlatMap(enabledPlugins, func(plugin plugin.Plugin) []plugin.Command {
			return plugin.Commands()
		}),
	}, nil
}

func (r *resolver) resolveNames() (map[string]ref.Ref, error) {
	for _, name := range r.disabledPlugins {
		if _, ok := r.plugins[name.String()]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownPlugin, name.String())
		}
	}

	disabledNamesMap := namesToMap(r.disabledPlugins)
	namesMap := namesToMap(r.enabledPlugins)

	resultMap := make(map[string]ref.Ref)

	var addPlugin func(name ref.Ref) error
	addPlugin = func(name ref.Ref) error {
		if _, ok := r.plugins[name.String()]; !ok {
			return fmt.Errorf("%w: %s", ErrUnknownPlugin, name)
		}

		if _, ok := disabledNamesMap[name.String()]; ok {
			return fmt.Errorf("%w: %s", ErrDisabled, name)
		}

		if _, ok := resultMap[name.String()]; ok {
			return nil
		}

		for _, dep := range r.DependenciesOf(name) {
			if err := addPlugin(dep); err != nil {
				return fmt.Errorf("%w: %s: %w", ErrDependency, name, err)
			}
		}

		resultMap[name.String()] = name

		return nil
	}

	if r.defaultPlugins {
		for _, plugin := range r.plugins {
			if plugin.EnabledByDefault() {
				_ = addPlugin(plugin.Ref())
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
