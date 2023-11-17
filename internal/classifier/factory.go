package classifier

import (
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/database/persistence"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In
	Persistence  persistence.Persistence
	SubResolvers []SubResolver `group:"content_resolvers"`
	Dao          *dao.Query
	Logger       *zap.SugaredLogger
}

type Result struct {
	fx.Out
	Classifier   Classifier
	Duration     prometheus.Collector `group:"prometheus_collectors"`
	SuccessTotal prometheus.Collector `group:"prometheus_collectors"`
	NoMatchTotal prometheus.Collector `group:"prometheus_collectors"`
	ErrorTotal   prometheus.Collector `group:"prometheus_collectors"`
}

func New(p Params) Result {
	collector := newPrometheusCollectorResolver(resolver{
		subResolvers: p.SubResolvers,
		logger:       p.Logger.Named("content_classifier"),
	})
	return Result{
		Classifier: classifier{
			resolver:    collector,
			dao:         p.Dao,
			persistence: p.Persistence,
		},
		Duration:     collector.duration,
		SuccessTotal: collector.successTotal,
		NoMatchTotal: collector.noMatchTotal,
		ErrorTotal:   collector.errorTotal,
	}
}
