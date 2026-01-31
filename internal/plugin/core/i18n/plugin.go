package i18n

import (
	"github.com/bitmagnet-io/bitmagnet/i18n"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
	"go.uber.org/fx"
)

type deps struct{}

var (
	Ref = ref.Root.MustSub("i18n")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithDescription[deps]("Provides translations for backend services"),
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithFxOption[deps](
			fx.Supply(i18n.Bundle),
		),
	)
)
