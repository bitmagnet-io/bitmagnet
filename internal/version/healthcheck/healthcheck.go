package healthcheck

import (
	"github.com/bitmagnet-io/bitmagnet/internal/version"
	"github.com/hellofresh/health-go/v5"
	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	HealthcheckOption health.Option `group:"healthcheck_options"`
}

func New() Result {
	return Result{
		HealthcheckOption: health.WithComponent(health.Component{
			Name:    "bitmagnet",
			Version: version.GitTag,
		}),
	}
}
