package server

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/prometheus/client_golang/prometheus"
	"net/netip"
	"time"
)

type prometheusCollector struct {
	server            Server
	queryDuration     *prometheus.HistogramVec
	querySuccessTotal *prometheus.CounterVec
	queryErrorTotal   *prometheus.CounterVec
	queryConcurrency  *prometheus.GaugeVec
}

const labelQuery = "query"

var labelNames = []string{labelQuery}

func newPrometheusCollector(server Server) prometheusCollector {
	return prometheusCollector{
		server: server,
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

func (l prometheusCollector) Ready() <-chan struct{} {
	return l.server.Ready()
}

func (l prometheusCollector) Query(ctx context.Context, addr netip.AddrPort, q string, args dht.MsgArgs) (dht.RecvMsg, error) {
	labels := prometheus.Labels{labelQuery: q}
	l.queryConcurrency.With(labels).Inc()
	start := time.Now()
	res, err := l.server.Query(ctx, addr, q, args)
	l.queryConcurrency.With(labels).Dec()
	if err == nil {
		l.queryDuration.With(labels).Observe(time.Since(start).Seconds())
		l.querySuccessTotal.With(labels).Inc()
	} else {
		l.queryErrorTotal.With(labels).Inc()
	}
	return res, err
}
