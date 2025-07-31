package registry

import (
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
)

type Option func(*Registry)

func WithWorker(name string, run runner.Provider, options ...worker.Option) Option {
	return func(r *Registry) {
		r.workers[name] = worker.NewWorker(
			run,
			append(options, worker.WithLogger(r.logger.Named(name)))...,
		)
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
