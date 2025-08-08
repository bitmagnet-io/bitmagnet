package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/hashicorp/go-metrics"
)

// Registry tracks all metrics within the app.
type Registry struct {
	mtx         *sync.RWMutex
	serviceName string
	metrics     *metrics.Metrics
	stopped     chan struct{}
	labelValues []LabelValue
	sink        *sink
	inmem       *metrics.InmemSink
	startTime   time.Time
	components  map[string]struct{}
}

type registryBuilder struct {
	serviceName       string
	labelValues       []LabelValue
	sinks             metrics.FanoutSink
	inmemSinkInterval time.Duration
	inmemSinkRetain   time.Duration
}

// NewRegistry creates a new metrics registry with the provided options.
func NewRegistry(serviceName string, options ...Option) (*Registry, error) {
	builder := &registryBuilder{
		serviceName:       serviceName,
		inmemSinkInterval: time.Minute,
		inmemSinkRetain:   time.Minute * 60 * 24 * 30,
	}

	for _, option := range options {
		if err := option(builder); err != nil {
			return nil, err
		}
	}

	var mtx sync.RWMutex

	snk := &sink{
		mtx:         &mtx,
		serviceName: builder.serviceName,
		totals:      make(Totals),
		averages:    make(Averages),
	}

	libCfg := metrics.DefaultConfig(builder.serviceName)

	libCfg.EnableHostname = false
	libCfg.EnableRuntimeMetrics = false
	libCfg.EnableTypePrefix = false

	inmem := metrics.NewInmemSink(builder.inmemSinkInterval, builder.inmemSinkRetain)

	libMetrics, err := metrics.New(libCfg, append(metrics.FanoutSink{
		snk,
		inmem,
	}, builder.sinks...))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics registry: %w", err)
	}

	registry := &Registry{
		mtx:         &mtx,
		serviceName: builder.serviceName,
		metrics:     libMetrics,
		stopped:     make(chan struct{}),
		labelValues: builder.labelValues,
		sink:        snk,
		inmem:       inmem,
		startTime:   time.Now(),
		components:  make(map[string]struct{}),
	}

	return registry, nil
}

func (r *Registry) NewComponent(ref ref.Ref) (*Component, error) {
	if _, exists := r.components[ref.String()]; exists {
		return nil, fmt.Errorf("component %s already exists", ref.String())
	}

	return &Component{
		ref:      ref,
		registry: r,
	}, nil
}

func (r *Registry) MustNewComponent(ref ref.Ref) *Component {
	component, err := r.NewComponent(ref)
	if err != nil {
		panic(err)
	}

	return component
}

// Totals returns all the accumulated total metrics.
func (r *Registry) Totals() Totals {
	return r.sink.Totals()
}

// Averages returns all the accumulated average metrics.
func (r *Registry) Averages() Averages {
	return r.sink.Averages()
}

// StartTime returns the time when this registry was initialized.
func (r *Registry) StartTime() time.Time {
	return r.startTime
}

// ExecutionTime returns the duration since this registry was initialized.
func (r *Registry) ExecutionTime() time.Duration {
	return time.Since(r.startTime)
}

// Close shuts down the registry and stops collecting metrics.
func (r *Registry) Close() {
	close(r.stopped)
	r.metrics.Shutdown()
}
