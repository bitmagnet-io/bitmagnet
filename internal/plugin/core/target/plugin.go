package target

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	plugin_search "github.com/bitmagnet-io/bitmagnet/internal/plugin/core/search"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/search"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/target"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type deps struct {
	fx.In
}

var (
	Ref = ref.Root.MustSub("target")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Allows torrents to be sent to external targets"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithDependencies[deps](
			plugin_search.Ref,
		),
		builder.WithFxOption[deps](
			fx.Provide(
				fx.Annotate(
					func(targets [][]target.TorrentContentTarget, search search.TorrentContent) target.Registry {
						return target.NewRegistry(
							slice.FlatMap(
								targets,
								func(targets []target.TorrentContentTarget) []target.TorrentContentTarget {
									return targets
								},
							),
							search,
						)
					},
					fx.ParamTags(`group:"torrent_targets"`),
				),
			),
		),
	)
)
