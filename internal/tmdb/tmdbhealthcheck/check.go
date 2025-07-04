package tmdbhealthcheck

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/tmdb"
)

func New(
	enabled bool,
	client tmdb.Client,
) health.CheckerOption {
	return health.WithPeriodicCheck(
		time.Minute*5,
		time.Second*5,
		health.Check{
			Name:    "tmdb",
			Timeout: time.Second * 30,
			IsActive: func() bool {
				return enabled
			},
			Check: func(ctx context.Context) error {
				return client.ValidateAPIKey(ctx)
			},
		},
	)
}
