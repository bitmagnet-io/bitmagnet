package healthcheck

import (
	"context"
	"fmt"
	"github.com/hellofresh/health-go/v5"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"time"
)

type Params struct {
	fx.In
	GormDb *gorm.DB
}

type Result struct {
	fx.Out
	Config health.Option `group:"healthcheck_options"`
}

func New(p Params) Result {
	return Result{
		Config: health.WithChecks(health.Config{
			Name:    "postgres",
			Timeout: time.Second * 5,
			Check: func(ctx context.Context) error {
				db, dbErr := p.GormDb.DB()
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
