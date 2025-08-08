package metrics

import (
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics/torrentmetrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"go.uber.org/fx"
)

type (
	config struct{}
	deps   struct{}
)

var (
	Ref = core.Ref.MustSub("metrics")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithFxOption[config, deps](
			fx.Provide(
				fx.Annotate(
					func(options []metrics.Option) (*metrics.Registry, error) {
						return metrics.NewRegistry("bitmagnet", options...)
					},
					fx.ParamTags(`group:"metrics_options"`),
				),
				// todo: Move this
				torrentmetrics.New,
			),
		),
	)
)
