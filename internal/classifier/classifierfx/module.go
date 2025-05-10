package classifierfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"workflow",
		configfx.NewConfigModule[classifier.Config]("classifier", classifier.NewDefaultConfig()),
		fx.Provide(
			classifier.New,
		),
	)
}
