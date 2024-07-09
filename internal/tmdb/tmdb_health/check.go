package tmdb_health

import (
	"context"
	"github.com/bitmagnet-io/bitmagnet/internal/boilerplate/lazy"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
	"time"
)

func NewCheck(
	enabled bool,
	client lazy.Lazy[tmdb.Client],
) health.Check {
	return health.Check{
		Name:    "tmdb",
		Timeout: time.Second * 30,
		IsActive: func() bool {
			return enabled
		},
		Check: func(ctx context.Context) error {
			c, err := client.Get()
			if err != nil {
				return err
			}
			return c.ValidateApiKey(ctx)
		},
	}
}
