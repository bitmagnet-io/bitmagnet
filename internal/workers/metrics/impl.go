package metrics

import (
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
)

const (
	NameBuffered     = "buffered"
	NameDiscarded    = "discarded"
	NameDeduplicated = "deduplicated"
	NameDequeued     = "dequeued"
	NameFlushing     = "flushing"
	NameFlushed      = "flushed"
)

func New(component *metrics.Component, labels ...metrics.LabelValue) (Adapter, error) {
	bufferedGauge, err := component.NewGauge(NameBuffered, labels...)
	if err != nil {
		return nil, err
	}

	discardedCounter, err := component.NewCounter(NameDiscarded, labels...)
	if err != nil {
		return nil, err
	}

	deduplicatorCounter, err := component.NewCounter(NameDeduplicated, labels...)
	if err != nil {
		return nil, err
	}

	dequeuedSampler, err := component.NewSampler(NameDequeued, labels...)
	if err != nil {
		return nil, err
	}

	flushingGauge, err := component.NewGauge(NameFlushing, labels...)
	if err != nil {
		return nil, err
	}

	flushedSampler, err := component.NewSampler(NameFlushed, labels...)
	if err != nil {
		return nil, err
	}

	return &adapter{
		bufferedGauge:       bufferedGauge,
		discardedCounter:    discardedCounter,
		deduplicatedCounter: deduplicatorCounter,
		dequeuedSampler:     dequeuedSampler,
		flushingGauge:       flushingGauge,
		flushedSampler:      flushedSampler,
	}, nil
}

func MustNew(component *metrics.Component, labels ...metrics.LabelValue) Adapter {
	adapter, err := New(component, labels...)
	if err != nil {
		panic(err)
	}

	return adapter
}

type adapter struct {
	mtx                 sync.Mutex
	bufferedCount       int
	flushingCount       int
	bufferedGauge       *metrics.Gauge
	discardedCounter    *metrics.Counter
	deduplicatedCounter *metrics.Counter
	dequeuedSampler     *metrics.Sampler
	flushingGauge       *metrics.Gauge
	flushedSampler      *metrics.Sampler
}

func (a *adapter) IncrAdded(n int) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.bufferedCount += n
	a.bufferedGauge.Set(a.bufferedCount)
}

func (a *adapter) IncrDiscarded() {
	a.discardedCounter.Incr()
}

func (a *adapter) IncrDeduplicated() {
	a.deduplicatedCounter.Incr()
}

func (a *adapter) IncrDequeued(latency time.Duration) {
	a.dequeuedSampler.Add(float32(latency.Seconds()))

	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.bufferedCount--
	a.flushingCount++
	a.bufferedGauge.Set(a.bufferedCount)
	a.flushingGauge.Set(a.flushingCount)
}

func (a *adapter) IncrFlushed(latency time.Duration) {
	a.flushedSampler.Add(float32(latency.Seconds()))

	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.flushingCount--
	a.flushingGauge.Set(a.flushingCount)
}

func (a *adapter) Reset() {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.bufferedCount = 0
	a.flushingCount = 0
	a.bufferedGauge.Set(0)
	a.flushingGauge.Set(0)
}
