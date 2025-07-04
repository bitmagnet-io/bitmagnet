package validationfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/validation"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"validation",
		fx.Provide(
			fx.Annotate(
				validation.New,
				fx.ParamTags(`group:"validator_options"`),
			),
		),
	)
}
