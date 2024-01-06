package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Search       lazy.Lazy[search.Search]
	SubResolvers []lazy.Lazy[SubResolver] `group:"content_resolvers"`
	Dao          lazy.Lazy[*dao.Query]
	Logger       *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Classifier   lazy.Lazy[Classifier]
	Duration     prometheus.Collector `group:"prometheus_collectors"`
	SuccessTotal prometheus.Collector `group:"prometheus_collectors"`
	NoMatchTotal prometheus.Collector `group:"prometheus_collectors"`
	ErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
}

func New(p Params) Result {
	collector := newPrometheusCollector()
	return Result{
		Classifier: lazy.New(func() (Classifier, error) {
			s, err := p.Search.Get()
			if err != nil {
				return classifier{}, err
			}
			d, err := p.Dao.Get()
			if err != nil {
				return classifier{}, err
			}
			return classifier{
				resolver: prometheusCollectorResolver{
					prometheusCollector: collector,
				},
				dao:    d,
				search: s,
			}, nil
		}),
		Duration:     collector.duration,
		SuccessTotal: collector.successTotal,
		NoMatchTotal: collector.noMatchTotal,
		ErrorTotal:   collector.errorTotal,
	}
}
