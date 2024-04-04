package classifierfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/classifier"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"workflow",
		fx.Provide(
			classifier.New,
		),
	)
}
