package metainforequester

import (
	"net"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/concurrency"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

type Params struct {
	fx.In
	Config Config
	Logger *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Requester           Requester
	RequestDuration     prometheus.Collector `group:"prometheus_collectors"`
	RequestSuccessTotal prometheus.Collector `group:"prometheus_collectors"`
	RequestErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
	RequestConcurrency  prometheus.Collector `group:"prometheus_collectors"`
}

func New(p Params) Result {
	collector := newPrometheusCollector(requester{
		clientID: protocol.RandomPeerID(),
		timeout:  p.Config.RequestTimeout,
		dialer: &net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: -1,
		},
	})

	return Result{
		Requester: requestLimiter{
			requester: requestLogger{
				requester: collector,
				// we make way to many requests to usefully log everything, but having a sample is
				// helpful:
				logger: p.Logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
					return zapcore.NewSamplerWithOptions(core, time.Minute, 10, 0)
				})).Named("meta_info_requester"),
			},
			limiter: concurrency.NewKeyedLimiter(rate.Every(time.Second/2), 4, 1000, time.Second*20),
		},
		RequestDuration:     collector.requestDuration,
		RequestSuccessTotal: collector.requestSuccessTotal,
		RequestErrorTotal:   collector.requestErrorTotal,
		RequestConcurrency:  collector.requestConcurrency,
	}
}
