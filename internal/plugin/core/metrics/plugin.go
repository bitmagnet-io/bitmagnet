package metrics

import (
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
				// todo: Move this
				torrentmetrics.New,
			),
		),
	)
)
