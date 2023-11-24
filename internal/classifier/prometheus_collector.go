package classifier

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type prometheusCollectorResolver struct {
	resolver     Resolver
	duration     *prometheus.HistogramVec
	successTotal *prometheus.CounterVec
	noMatchTotal prometheus.Counter
	errorTotal   prometheus.Counter
}

const namespace = "bitmagnet"
const subsystem = "classifier"

const labelContentType = "content_type"
const labelContentSource = "content_source"

func newPrometheusCollectorResolver(resolver Resolver) prometheusCollectorResolver {
	return prometheusCollectorResolver{
		resolver: resolver,
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "duration_seconds",
			Help:      "A histogram of successful classification durations in seconds.",
			Buckets:   prometheus.ExponentialBuckets(0.1, 1.5, 5),
		}, []string{labelContentType, labelContentSource}),
		successTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "success_total",
			Help:      "A counter of successful classifications.",
		}, []string{labelContentType, labelContentSource}),
		noMatchTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "no_match_total",
			Help:      "A counter of classifications with no match.",
		}),
		errorTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "error_total",
			Help:      "A counter of failed classifications.",
		}),
	}
}

func (r prometheusCollectorResolver) Resolve(ctx context.Context, content model.TorrentContent) (model.TorrentContent, error) {
	start := time.Now()
	result, err := r.resolver.Resolve(ctx, content)
	contentType := "unknown"
	contentSource := "unknown"
	if err == nil {
		if result.ContentType.Valid {
			contentType = result.ContentType.ContentType.String()
		}
		if result.ContentSource.Valid {
			contentSource = result.ContentSource.String
		}
		r.successTotal.With(prometheus.Labels{
			labelContentType:   contentType,
			labelContentSource: contentSource,
		}).Inc()
	} else if errors.Is(err, ErrNoMatch) {
		r.noMatchTotal.Inc()
	} else {
		r.errorTotal.Inc()
		goto END
	}
	r.duration.With(prometheus.Labels{
		labelContentType:   contentType,
		labelContentSource: contentSource,
	}).Observe(time.Since(start).Seconds())
END:
	return result, err
}
