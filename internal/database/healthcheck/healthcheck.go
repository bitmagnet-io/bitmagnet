package healthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
)

func New(name string, p database.Provider) health.CheckerOption {
	return health.WithPeriodicCheck(
		time.Second*30,
		time.Second*1,
		health.Check{
			Name:    name,
			Timeout: time.Second * 5,
			IsActive: func() bool {
				return p.IsActive()
			},
			Check: func(ctx context.Context) error {
				pool, err := p.Pool()
				if err != nil {
					return fmt.Errorf("failed to get database connection: %w", err)
				}

				err = pool.Ping(ctx)
				if err != nil {
					return fmt.Errorf("failed to ping database: %w", err)
				}

				return nil
			},
		},
	)
}
