package boilerplatefx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/cli/clifx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/logging/loggingfx"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/validation/validationfx"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"boilerplate",
		clifx.New(),
		configfx.New(),
		loggingfx.New(),
		validationfx.New(),
	)
}
