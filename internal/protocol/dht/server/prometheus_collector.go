package server

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/prometheus/client_golang/prometheus"
	"net/netip"
	"time"
)

type prometheusCollector struct {
	queryDuration     *prometheus.HistogramVec
	querySuccessTotal *prometheus.CounterVec
	queryErrorTotal   *prometheus.CounterVec
	queryConcurrency  *prometheus.GaugeVec
}

const labelQuery = "query"

var labelNames = []string{labelQuery}

func newPrometheusCollector() prometheusCollector {
	return prometheusCollector{
		queryDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "query_duration_seconds",
			Help:      "A histogram of successful DHT query durations in seconds.",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		}, labelNames),
		querySuccessTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "query_success_total",
			Help:      "A counter of successful DHT queries.",
		}, labelNames),
		queryErrorTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "query_error_total",
			Help:      "A counter of failed DHT queries.",
		}, labelNames),
		queryConcurrency: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "query_concurrency",
			Help:      "Number of concurrent DHT queries.",
		}, labelNames),
	}
}

type prometheusServerWrapper struct {
	prometheusCollector
	server Server
}

func (s prometheusServerWrapper) start() error {
	return s.server.start()
}

func (s prometheusServerWrapper) stop() {
	s.server.stop()
}

func (s prometheusServerWrapper) Query(ctx context.Context, addr netip.AddrPort, q string, args dht.MsgArgs) (dht.RecvMsg, error) {
	labels := prometheus.Labels{labelQuery: q}
	s.queryConcurrency.With(labels).Inc()
	start := time.Now()
	res, err := s.server.Query(ctx, addr, q, args)
	s.queryConcurrency.With(labels).Dec()
	if err == nil {
		s.queryDuration.With(labels).Observe(time.Since(start).Seconds())
		s.querySuccessTotal.With(labels).Inc()
	} else {
		s.queryErrorTotal.With(labels).Inc()
	}
	return res, err
}
