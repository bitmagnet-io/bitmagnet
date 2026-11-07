package matchlabhealthcheck

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/matchlab"
)

// Readiness, not liveness: the MatchLab matcher must have loaded its index.
func NewCheck(cfg matchlab.Config, client lazy.Lazy[matchlab.Client]) health.Check {
	return health.Check{
		Name:     "matchlab",
		Timeout:  10 * time.Second,
		IsActive: cfg.Enabled,
		Check: func(ctx context.Context) error {
			c, err := client.Get()
			if err != nil {
				return err
			}
			return c.Health(ctx)
		},
	}
}
