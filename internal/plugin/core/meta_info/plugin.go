package meta_info

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/banning"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"go.uber.org/fx"
)

type (
	Config = metainforequester.Config

	deps struct{}
)

var (
	Ref = core.Ref.MustSub("meta_info")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithDependencies[Config, deps](
			config.Ref,
		),
		builder.WithDefaultConfig[Config, deps](metainforequester.NewDefaultConfig()),
		builder.WithFxOption[Config, deps](
			fx.Provide(
				metainforequester.New,
				banning.New,
			),
		),
	)
)
