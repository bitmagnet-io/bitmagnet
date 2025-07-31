package validation

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/validation"
	"go.uber.org/fx"
)

type (
	config struct{}
	deps   struct{}
)

var (
	Ref = core.Ref.MustSub("validation")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithFxOption[config, deps](
			fx.Provide(
				fx.Annotate(
					validation.New,
					fx.ParamTags(`group:"validator_options"`),
				),
			),
		),
	)
)
