package musicfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/music/discogs"
	// "github.com/bitmagnet-io/bitmagnet/internal/classifier/resolver"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"music",
		configfx.NewConfigModule[discogs.Config]("discogs", discogs.NewDefaultConfig()),
		fx.Provide(
			discogs.New,
			// resolver.New, this is already provide by the first resolver which is video
		),
	)
}
