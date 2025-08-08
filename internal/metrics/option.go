package metrics

import "github.com/hashicorp/go-metrics"

// Option is a function that configures a registryBuilder.
type Option func(*registryBuilder) error

// WithSink adds a metrics sink to the registry.
func WithSink(sink metrics.MetricSink) Option {
	return func(r *registryBuilder) error {
		r.sinks = append(r.sinks, sink)

		return nil
	}
}

// WithLabelsValues adds label values to the registry.
func WithLabelsValues(labelValues ...LabelValue) Option {
	return func(r *registryBuilder) error {
		r.labelValues = append(r.labelValues, labelValues...)

		return nil
	}
}
