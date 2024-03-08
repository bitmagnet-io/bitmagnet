package prometheus

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Collectors []prometheus.Collector `group:"prometheus_collectors"`
}

type Result struct {
	fx.Out
	Registry lazy.Lazy[*prometheus.Registry]
}

const Namespace = "bitmagnet"

func New(p Params) (Result, error) {
	return Result{
		Registry: lazy.New[*prometheus.Registry](func() (*prometheus.Registry, error) {
			registry := prometheus.NewRegistry()
			cs := append(
				[]prometheus.Collector{
					collectors.NewGoCollector(),
					collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
						Namespace: Namespace,
					}),
				},
				p.Collectors...,
			)
			for _, c := range cs {
				if err := registry.Register(c); err != nil {
					return nil, err
				}
			}
			return registry, nil
		}),
	}, nil
}
