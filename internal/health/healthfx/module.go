package healthfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"health",
		fx.Provide(
			health.New,
		),
	)
}
