package torznab

import (
	internalsearch "github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/adapter"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/httpserver"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// todo: Profiles!
type deps struct {
	fx.In
	Client torznab.Client
}

var (
	Ref = core.Ref.MustSub("torznab")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[deps](),
		builder.WithDependencies[deps](
			config.Ref,
			search.Ref,
		),
		// builder.WithDefaultConfig[deps](torznab.NewDefaultConfig()),
		builder.WithFxOption[deps](
			fx.Provide(
				func(s internalsearch.Search) torznab.Client {
					return adapter.New(s)
				},
			),
		),
		builder.WithGinOption(
			Ref,
			func(deps deps) gin.OptionFunc {
				return httpserver.Option(torznab.Config{}, deps.Client)
			},
		),
	)
)
