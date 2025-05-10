package healthcheck

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/lazy"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	DB lazy.Lazy[*sql.DB]
}

type Result struct {
	fx.Out
	Option health.CheckerOption `group:"health_check_options"`
}

func New(p Params) Result {
	return Result{
		Option: health.WithPeriodicCheck(
			time.Second*30,
			time.Second*1,
			health.Check{
				Name:    "postgres",
				Timeout: time.Second * 5,
				Check: func(ctx context.Context) error {
					db, dbErr := p.DB.Get()
					if dbErr != nil {
						return fmt.Errorf("failed to get database connection: %w", dbErr)
					}
					pingErr := db.PingContext(ctx)
					if pingErr != nil {
						return fmt.Errorf("failed to ping database: %w", pingErr)
					}
					return nil
				},
			}),
	}
}
