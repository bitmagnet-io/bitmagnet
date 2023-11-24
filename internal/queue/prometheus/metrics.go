package prometheus

import (
	"github.com/bitmagnet-io/bitmagnet/internal/queue/redis"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Redis *redis.Client
}

type Result struct {
	fx.Out
	Collector prometheus.Collector `group:"prometheus_collectors"`
}

func New(p Params) Result {
	return Result{
		Collector: metrics.NewQueueMetricsCollector(
			asynq.NewInspector(redis.Wrapper{Redis: p.Redis}),
		),
	}
}
