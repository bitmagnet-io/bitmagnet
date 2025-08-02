package worker

import "go.uber.org/zap"

type Option func(*Worker)

func WithLogger(logger *zap.Logger) Option {
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
