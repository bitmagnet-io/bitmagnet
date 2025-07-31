package worker

import "github.com/bitmagnet-io/bitmagnet/internal/logging"

type Option func(*Worker)

func WithLogger(logger logging.Logger) Option {
	return func(w *Worker) {
		w.logger = logger
	}
}

func WithDependencies(keys ...string) Option {
	return func(w *Worker) {
		for _, key := range keys {
			w.dependsOn[key] = struct{}{}
		}
	}
}

func WithAutostart() Option {
	return func(w *Worker) {
		w.autostart = true
	}
}
