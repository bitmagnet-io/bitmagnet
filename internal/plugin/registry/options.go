package registry

import "github.com/bitmagnet-io/bitmagnet/internal/ref"

type Option func(*resolver)

func WithDefaultPlugins() Option {
	return func(r *resolver) {
		r.defaultPlugins = true
	}
}

func WithEnabledPlugins(refs ...ref.Ref) Option {
	return func(r *resolver) {
		r.enabledPlugins = append(r.enabledPlugins, refs...)
	}
}

func WithDisabledPlugins(refs ...ref.Ref) Option {
	return func(r *resolver) {
		r.disabledPlugins = append(r.disabledPlugins, refs...)
	}
}
