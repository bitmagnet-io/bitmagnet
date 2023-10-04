package boilerplateappfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/app/cmd/config"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/app/cmd/worker"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/boilerplatefx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/cli/hooks"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/worker/workerfx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"app_boilerplate",
		boilerplatefx.New(),
		workerfx.New(),
		fx.Provide(
			hooks.New,
			configcmd.New,
			workercmd.New,
		),
	)
}
