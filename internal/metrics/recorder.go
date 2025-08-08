package metrics

import (
	"github.com/hashicorp/go-metrics"
)

// recorder is the base type for all metrics types.
type recorder struct {
	registry    *Registry
	ref         Ref
	labelValues []LabelValue
}

func (r recorder) Type() Type {
	return r.ref.Type
}

// Ref returns the Ref of the metric being recorded.
func (r recorder) Ref() Ref {
	return r.ref
}

func (r recorder) metricsLabels(labelValues []LabelValue) []metrics.Label {
	mlvs := metricsLabels(r.registry.labelValues)
	mlvs = append(mlvs, metricsLabels(r.labelValues)...)
	mlvs = append(mlvs, metricsLabels(labelValues)...)

	return mlvs
}
