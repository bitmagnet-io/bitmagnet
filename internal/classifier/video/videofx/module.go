package videofx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/tmdb"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"movie",
		configfx.NewConfigModule[tmdb.Config]("tmdb", tmdb.NewDefaultConfig()),
		fx.Provide(
			tmdb.New,
			video.New,
		),
	)
}
