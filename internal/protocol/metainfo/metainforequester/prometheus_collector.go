package metainforequester

import (
	"context"
	"net/netip"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol"
	"github.com/prometheus/client_golang/prometheus"
)

type prometheusCollector struct {
	requester           Requester
	requestDuration     prometheus.Histogram
	requestSuccessTotal prometheus.Counter
	requestErrorTotal   prometheus.Counter
	requestConcurrency  prometheus.Gauge
}

const (
	namespace = "bitmagnet"
	subsystem = "meta_info_requester"
)

func newPrometheusCollector(requester Requester) *prometheusCollector {
	return &prometheusCollector{
		requester: requester,
		requestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "duration_seconds",
			Help:      "Duration of successful meta info requests in seconds.",
			Buckets:   prometheus.DefBuckets,
		}),
		requestSuccessTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "success_total",
			Help:      "Total number of successful meta info requests.",
		}),
		requestErrorTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "error_total",
			Help:      "Total number of failed meta info requests.",
		}),
		requestConcurrency: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "concurrency",
			Help:      "Number of concurrent meta info requests.",
		}),
	}
}

func (l prometheusCollector) Request(ctx context.Context, infoHash protocol.ID, addr netip.AddrPort) (Response, error) {
	l.requestConcurrency.Inc()

	start := time.Now()
	resp, err := l.requester.Request(ctx, infoHash, addr)
	l.requestConcurrency.Dec()

	if err == nil {
		l.requestDuration.Observe(time.Since(start).Seconds())
		l.requestSuccessTotal.Inc()
	} else {
		l.requestErrorTotal.Inc()
	}

	return resp, err
}
