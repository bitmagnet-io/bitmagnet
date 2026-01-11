package dev

import (
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/builder"
	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/pkg/plugin"
)

type deps struct{}

var (
	Ref = ref.Root.MustSub("dev")

	Plugin = builder.NewPlugin(
		Ref,
		builder.WithActivation[deps](plugin.ActivationAlways),
		builder.WithCliCommand[deps](NewDevCommand()),
	)
)
