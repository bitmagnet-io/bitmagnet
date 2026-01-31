package metrics

import "github.com/bitmagnet-io/bitmagnet/internal/ref"

type Component struct {
	ref      ref.Ref
	registry *Registry
}

func (c *Component) Sub(name string) (*Component, error) {
	ref, err := c.ref.Sub(name)
	if err != nil {
		return nil, err
	}

	return c.registry.NewComponent(ref)
}

func (c *Component) MustSub(name string) *Component {
	component, err := c.Sub(name)
	if err != nil {
		panic(err)
	}

	return component
}

func (c *Component) NewCounter(name string, labels ...LabelValue) (*Counter, error) {
	ref, err := c.ref.Sub(name)
	if err != nil {
		return nil, err
	}

	return &Counter{
		recorder: recorder{
			registry: c.registry,
			ref: Ref{
				Ref:  ref,
				Type: TypeCounter,
			},
			labelValues: labels,
		},
	}, nil
}

func (c *Component) MustNewCounter(name string, labels ...LabelValue) *Counter {
	counter, err := c.NewCounter(name, labels...)
	if err != nil {
		panic(err)
	}

	return counter
}

func (c *Component) NewGauge(name string, labels ...LabelValue) (*Gauge, error) {
	ref, err := c.ref.Sub(name)
	if err != nil {
		return nil, err
	}

	return &Gauge{
		recorder: recorder{
			registry: c.registry,
			ref: Ref{
				Ref:  ref,
				Type: TypeGauge,
			},
			labelValues: labels,
		},
	}, nil
}

func (c *Component) MustNewGauge(name string, labels ...LabelValue) *Gauge {
	gauge, err := c.NewGauge(name, labels...)
	if err != nil {
		panic(err)
	}

	return gauge
}

func (c *Component) NewSampler(name string, labels ...LabelValue) (*Sampler, error) {
	ref, err := c.ref.Sub(name)
	if err != nil {
		return nil, err
	}

	return &Sampler{
		recorder: recorder{
			registry: c.registry,
			ref: Ref{
				Ref:  ref,
				Type: TypeSampler,
			},
			labelValues: labels,
		},
	}, nil
}

func (c *Component) MustNewSampler(name string, labels ...LabelValue) *Sampler {
	sampler, err := c.NewSampler(name, labels...)
	if err != nil {
		panic(err)
	}

	return sampler
}
