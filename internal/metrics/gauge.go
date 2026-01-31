package metrics

// Gauge represents a metric that can be set to a specific value
type Gauge struct {
	recorder
}

// Set sets the value of the gauge for the provided labels.
func (g *Gauge) Set(value int, labels ...LabelValue) {
	g.registry.metrics.SetGaugeWithLabels(g.ref.Path(), float32(value), g.metricsLabels(labels))
}

// Value returns the value of the gauge for the provided labels.
// The labels must be an exact match of those used when the value was Set.
func (g *Gauge) Value(labels ...LabelValue) int {
	return g.registry.Totals().Value(g.ref, labels...)
}

// ValuesByLabel returns the gauge values grouped by the specified label
func (g *Gauge) ValuesByLabel(label Label) map[string]int {
	return g.registry.Totals().ValuesByLabel(g.ref, label)
}
