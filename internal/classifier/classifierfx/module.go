package classifierfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/extension"
	"github.com/bitmagnet-io/bitmagnet/internal/classifier/video/videofx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"classifier",
		fx.Provide(
			classifier.New,
			extension.New,
		),
		videofx.New(),
	)
}
