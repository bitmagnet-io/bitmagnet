package workflowfx

import (
	"github.com/bitmagnet-io/bitmagnet/internal/workflow"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Module(
		"workflow",
		fx.Provide(
			workflow.New,
		),
	)
}
