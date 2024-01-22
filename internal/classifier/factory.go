package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sort"
)

type Params struct {
	fx.In
	SubResolvers []lazy.Lazy[SubClassifier] `group:"content_resolvers"`
	Logger       *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Classifier lazy.Lazy[Classifier]
	//Duration     prometheus.Collector `group:"prometheus_collectors"`
	//SuccessTotal prometheus.Collector `group:"prometheus_collectors"`
	//NoMatchTotal prometheus.Collector `group:"prometheus_collectors"`
	//ErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
}

func New(p Params) Result {
	//collector := newPrometheusCollector()
	return Result{
		Classifier: lazy.New(func() (Classifier, error) {
			//s, err := p.Search.Get()
			//if err != nil {
			//	return nil, err
			//}
			//d, err := p.Dao.Get()
			//if err != nil {
			//	return nil, err
			//}
			subClassifiers := make([]SubClassifier, 0, len(p.SubResolvers)+1)
			for _, subResolver := range p.SubResolvers {
				r, err := subResolver.Get()
				if err != nil {
					return nil, err
				}
				subClassifiers = append(subClassifiers, r)
			}
			subClassifiers = append(subClassifiers, fallbackClassifier{})
			sort.Slice(subClassifiers, func(i, j int) bool {
				return subClassifiers[i].Priority() < subClassifiers[j].Priority()
			})
			return classifier{subClassifiers, p.Logger}, nil
		}),
		//Duration:     collector.duration,
		//SuccessTotal: collector.successTotal,
		//NoMatchTotal: collector.noMatchTotal,
		//ErrorTotal:   collector.errorTotal,
	}
}
