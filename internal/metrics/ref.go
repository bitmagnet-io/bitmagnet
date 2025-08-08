package metrics

import (
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
)

type Ref struct {
	ref.Ref
	Type
}

// Type represents the type of metric
type Type string

const (
	// TypeCounter represents a counter metric type that accumulates over time
	TypeCounter Type = "counter"
	// TypeGauge represents a gauge metric type that can be set to a specific value
	TypeGauge Type = "gauge"
	// TypeSampler represents a sampler metric type that collects samples
	TypeSampler Type = "sample"
)
