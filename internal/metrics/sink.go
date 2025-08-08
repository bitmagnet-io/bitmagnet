package metrics

import (
	"sync"

	"github.com/hashicorp/go-metrics"
)

// sink is the core metrics.MetricSink implementation, and tracks Totals and Averages for all metrics.
type sink struct {
	mtx         *sync.RWMutex
	serviceName string
	totals      Totals
	averages    Averages
}

var _ metrics.MetricSink = (*sink)(nil)

func (t *sink) EmitKey(_ []string, _ float32) {
	// unsupported
}

func (t *sink) AddSample(key []string, val float32) {
	t.AddSampleWithLabels(key, val, nil)
}

func (t *sink) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
	if key[0] != t.serviceName {
		return
	}

	hash := flattenRefLabels(key, labels...)

	t.mtx.Lock()

	t.totals[hash]++

	avg := t.averages[hash]
	avg.Count++
	avg.Sum += float64(val)
	t.averages[hash] = avg

	t.mtx.Unlock()
}

func (t *sink) IncrCounter(key []string, val float32) {
	t.IncrCounterWithLabels(key, val, nil)
}

func (t *sink) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
	if key[0] != t.serviceName {
		return
	}

	hash := flattenRefLabels(key, labels...)

	t.mtx.Lock()
	t.totals[hash] += int(val)
	t.mtx.Unlock()
}

func (t *sink) SetGauge(key []string, val float32) {
	t.SetGaugeWithLabels(key, val, nil)
}

func (t *sink) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
	if key[0] != t.serviceName {
		return
	}

	hash := flattenRefLabels(key, labels...)

	t.mtx.Lock()
	t.totals[hash] = int(val)
	t.mtx.Unlock()
}
