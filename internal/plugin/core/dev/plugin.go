package dev

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/core"
)

type deps struct{}

var (
	Ref = core.Ref.MustSub("dev")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithCliCommand[deps](NewDevCommand()),
	)
)
