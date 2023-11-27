package blockingfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/blocking"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"blocking",
		fx.Provide(
			blocking.New,
		),
	)
}
