package metainforequester

import (
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/bitmagnet-io/bitmagnet/internal/semaphore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

// type Params struct {
// 	fx.In
// 	Config Config
// 	Logger *zap.Logger
// }

// type Result struct {
// 	fx.Out
// 	Requester           Requester
// 	RequestDuration     prometheus.Collector `group:"prometheus_collectors"`
// 	RequestSuccessTotal prometheus.Collector `group:"prometheus_collectors"`
// 	RequestErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
// 	RequestConcurrency  prometheus.Collector `group:"prometheus_collectors"`
// }

func New(
	maxConcurrency *atomic.Value[MaxConcurrency],
	dialTimeout *atomic.Value[DialTimeout],
	requestTimeout *atomic.Value[RequestTimeout],
	metrics *metrics.Component,
	logger *zap.Logger,
) Requester {
	mappedConcurrency, _ := atomic.Map(maxConcurrency, func(value MaxConcurrency) int {
		return int(value)
	})

	sem, _ := semaphore.NewAtomic(mappedConcurrency)

	return &requestLimiter{
		requester: &requesterSemaphore{
			Requester: &requestLogger{
				requester: newRequesterMetrics(
					&requester{
						clientID:       protocol.RandomPeerID(),
						dialTimeout:    dialTimeout,
						requestTimeout: requestTimeout,
					},
					metrics,
				),
				// we make way to many requests to usefully log everything, but having a sample is
				// helpful:
				logger: logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
					return zapcore.NewSamplerWithOptions(core, time.Minute, 10, 0)
				})),
			},
			semaphore: sem,
		},
		limiter: concurrency.NewKeyedLimiter(rate.Every(time.Second/2), 4, 10_000, time.Minute),
	}
	// 	RequestDuration:     collector.requestDuration,
	// 	RequestSuccessTotal: collector.requestSuccessTotal,
	// 	RequestErrorTotal:   collector.requestErrorTotal,
	// 	RequestConcurrency:  collector.requestConcurrency,
	// }
}
