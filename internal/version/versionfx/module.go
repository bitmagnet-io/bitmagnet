package versionfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/version/healthcheck"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"version",
		fx.Provide(healthcheck.New),
	)
}
