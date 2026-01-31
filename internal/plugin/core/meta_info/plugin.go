package meta_info

import (
	"github.com/bitmagnet-io/bitmagnet/internal/atomic"
	"github.com/bitmagnet-io/bitmagnet/internal/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	plugin_metrics "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/metrics"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/banning"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type deps struct {
	fx.In
}

var (
	Ref = ref.Root.MustSub("meta_info")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides services for downloading torrent meta info from peers"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			config.Ref,
			plugin_metrics.Ref,
		),
		builder.WithConfig[deps](Ref.MustSub("max_concurrency"), metainforequester.ParamMaxConcurrency),
		builder.WithConfig[deps](Ref.MustSub("dial_timeout"), metainforequester.ParamDialTimeout),
		builder.WithConfig[deps](Ref.MustSub("request_timeout"), metainforequester.ParamRequestTimeout),
		builder.WithFxOption[deps](
			fx.Provide(
				func(
					maxConcurrency *atomic.Value[metainforequester.MaxConcurrency],
					dialTimeout *atomic.Value[metainforequester.DialTimeout],
					requestTimeout *atomic.Value[metainforequester.RequestTimeout],
					metrics *metrics.Registry,
					logger *zap.Logger,
				) metainforequester.Requester {
					return metainforequester.New(
						maxConcurrency,
						dialTimeout,
						requestTimeout,
						metrics.MustNewComponent(Ref),
						logger.Named(Ref.String()),
					)
				},
				banning.New,
			),
		),
	)
)
