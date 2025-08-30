package registry

import (
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
)

type Option func(*Registry)

func WithWorker(ref ref.Ref, run runner.Provider, options ...worker.Option) Option {
	return func(r *Registry) {
		r.workers.Set(ref, worker.NewWorker(
			ref,
			run,
			append(options, worker.WithLogger(r.logger.Named(ref.String())))...,
		))
	}
}

func OptionNop() Option {
	return func(r *Registry) {}
}

func Options(options ...Option) Option {
	return func(r *Registry) {
		for _, option := range options {
			option(r)
		}
	}
}
