package healthcheck

import (
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	HealthOption health.CheckerOption `group:"health_check_options"`
}

func New() Result {
	return Result{
		HealthOption: health.WithInfo(map[string]any{
			"name":    "bitmagnet",
			"version": version.GitTag,
		}),
	}
}
