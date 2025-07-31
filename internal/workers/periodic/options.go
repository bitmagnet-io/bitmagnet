package periodic

import "time"

type Option func(*worker)

func WithQuickShutdown() Option {
	return func(wrk *worker) {
		wrk.quickShutdown = true
	}
}

func WithInitialInterval(interval time.Duration) Option {
	return func(wrk *worker) {
		wrk.initialInterval = interval
	}
}
