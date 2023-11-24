package responder

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type prometheusCollector struct {
	responder         Responder
	queryDuration     *prometheus.HistogramVec
	querySuccessTotal *prometheus.CounterVec
	queryErrorTotal   *prometheus.CounterVec
	queryConcurrency  *prometheus.GaugeVec
}

const labelQuery = "query"

var labelNames = []string{labelQuery}

func newPrometheusCollector(responder Responder) prometheusCollector {
	return prometheusCollector{
		responder: responder,
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

func (l prometheusCollector) Respond(ctx context.Context, msg dht.RecvMsg) (dht.Return, error) {
	labels := prometheus.Labels{labelQuery: normalizeQuery(msg.Msg.Q)}
	l.queryConcurrency.With(labels).Inc()
	start := time.Now()
	ret, err := l.responder.Respond(ctx, msg)
	l.queryConcurrency.With(labels).Dec()
	if err == nil {
		l.queryDuration.With(labels).Observe(time.Since(start).Seconds())
		l.querySuccessTotal.With(labels).Inc()
	} else {
		l.queryErrorTotal.With(labels).Inc()
	}
	return ret, err
}

func normalizeQuery(q string) string {
	switch q {
	case dht.QPing:
		return dht.QPing
	case dht.QFindNode:
		return dht.QFindNode
	case dht.QGetPeers:
		return dht.QGetPeers
	case dht.QSampleInfohashes:
		return dht.QSampleInfohashes
	default:
		return "unknown"
	}
}
