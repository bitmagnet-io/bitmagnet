package metrics

// Counter represents a metric that accumulates values over time
type Counter struct {
	recorder
}

// Incr increments the counter by 1 with the specified labels
func (c *Counter) Incr(labels ...LabelValue) {
	c.IncrN(1, labels...)
}

// IncrN increments the counter by n with the specified labels
func (c *Counter) IncrN(n int, labels ...LabelValue) {
	if n == 0 {
		return
	}

	c.registry.metrics.IncrCounterWithLabels(c.Ref().Path(), float32(n), c.metricsLabels(labels))
}

// Total returns the current total value of the counter with the specified labels
func (c *Counter) Total(labels ...LabelValue) int {
	return c.registry.Totals().Value(c.ref, labels...)
}

// TotalsByLabel returns the counter totals grouped by the specified label
func (c *Counter) TotalsByLabel(label Label) map[string]int {
	return c.registry.Totals().ValuesByLabel(c.ref, label)
}
