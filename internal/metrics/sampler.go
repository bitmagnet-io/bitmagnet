package metrics

// Sampler collects samples of a metric over time for statistical analysis.
type Sampler struct {
	recorder
}

// Add adds a new sample value with optional labels.
func (s *Sampler) Add(value float32, labels ...LabelValue) {
	s.registry.metrics.AddSampleWithLabels(s.Ref().Path(), value, s.metricsLabels(labels))
}

// Average returns the average value of all samples collected.
func (s *Sampler) Average() Average {
	return s.registry.Averages().Average(s.ref)
}

// AveragesByLabel returns a map of average values grouped by a specific label.
func (s *Sampler) AveragesByLabel(label Label) map[string]Average {
	return s.registry.Averages().AveragesByLabel(s.ref, label)
}
