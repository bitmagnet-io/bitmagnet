package torznab

import (
	internalsearch "github.com/bitmagnet-io/bitmagnet/internal/database/search"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/config"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/search"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/adapter"
	"github.com/bitmagnet-io/bitmagnet/internal/torznab/httpserver"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// todo: Profiles!
type deps struct {
	fx.In
	Client torznab.Client
}

var (
	Ref = ref.Root.MustSub("torznab")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Runs the Torznab API server on the /torznab endpoint"),
		builder.WithActivation[deps](plugin.ActivationEnabled),
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
			0,
			func(deps deps) gin.OptionFunc {
				return httpserver.Option(torznab.Config{}, deps.Client)
			},
		),
	)
)
