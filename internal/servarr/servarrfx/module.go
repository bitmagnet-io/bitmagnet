package servarrfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/servarr"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"servarr",
		configfx.NewConfigModule[servarr.Config]("servarr", servarr.NewDefaultConfig()),
	)
}
