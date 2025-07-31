package config

import (
	"github.com/bitmagnet-io/bitmagnet/internal/config/configfx"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core/validation"
)

type (
	config struct{}

	deps struct{}
)

var (
	Ref = core.Ref.MustSub("config")

	Plugin = builder.CreatePlugin(
		Ref,
		builder.WithEnabledByDefault[config, deps](),
		builder.WithDependencies[config, deps](
			validation.Ref,
		),
		builder.WithFxOption[config, deps](configfx.New()),
	)
)
