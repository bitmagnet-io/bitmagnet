package worker

import (
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"go.uber.org/zap"
)

type Option func(*Worker)

func Options(options ...Option) Option {
	return func(w *Worker) {
		for _, option := range options {
			option(w)
		}
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(w *Worker) {
		w.logger = logger
	}
}

func WithDependencies(refs ...ref.Ref) Option {
	return func(w *Worker) {
		for _, ref := range refs {
			w.dependsOn.Set(ref, struct{}{})
		}
	}
}

func WithAutostart(autostart bool) Option {
	return func(w *Worker) {
		w.autostart = autostart
	}
}

func ShortLived() Option {
	return func(w *Worker) {
		w.shortLived = true
	}
}
