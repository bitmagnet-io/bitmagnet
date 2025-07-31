package bundle

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type Bundle interface {
	Ref() ref.Ref
	Plugins() []plugin.Plugin
}

type bundle struct {
	ref     ref.Ref
	plugins map[string]plugin.Plugin
}

func (b *bundle) Ref() ref.Ref {
	return b.ref
}

func (b *bundle) Plugins() []plugin.Plugin {
	plugins := slices.Collect(maps.Values(b.plugins))
	slices.SortFunc(plugins, func(a, b plugin.Plugin) int {
		return cmp.Compare(a.Ref().String(), b.Ref().String())
	})

	return plugins
}

var _ Bundle = (*bundle)(nil)

func MustNew(
	bundleRef ref.Ref,
	bundlePlugins ...plugin.Plugin,
) Bundle {
	ret, err := New(bundleRef, bundlePlugins...)
	if err != nil {
		panic(err)
	}

	return ret
}

func New(
	bundleRef ref.Ref,
	bundlePlugins ...plugin.Plugin,
) (Bundle, error) {
	if !bundleRef.IsDescendentOf(ref.Root) {
		return nil, fmt.Errorf("%w: %w: %s", Err, ErrInvalidRef, bundleRef)
	}

	pluginsMap := make(map[string]plugin.Plugin)

	var err error

	for _, plugin := range bundlePlugins {
		pluginRef := plugin.Ref()

		if !pluginRef.IsDescendentOf(bundleRef) {
			err = fmt.Errorf("%w: %s", ErrInvalidRef, pluginRef)
			break
		}

		if _, ok := pluginsMap[pluginRef.String()]; ok {
			err = fmt.Errorf("%w: %s", ErrAlreadyRegistered, pluginRef)
		}

		pluginsMap[pluginRef.String()] = plugin
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", Err, bundleRef, err)
	}

	return &bundle{
		ref:     bundleRef,
		plugins: pluginsMap,
	}, nil
}
