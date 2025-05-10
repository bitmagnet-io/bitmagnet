package metainfofx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/banning"
	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo/metainforequester"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"metainfo",
		configfx.NewConfigModule[metainforequester.Config](
			"metainfo_requester",
			metainforequester.NewDefaultConfig(),
		),
		fx.Provide(
			metainforequester.New,
			banning.New,
		),
	)
}
