package gormcmd

import (
	"context"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/app"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/gen"
	gormlogger "github.com/bitmagnet-io/bitmagnet/internal/database/logger"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormliblogger "gorm.io/gorm/logger"
)

var Command = &cli.Command{
	Name:  "gorm-gen",
	Usage: "Generates GORM models from the database schema",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fxApp := app.New(
			ctx,
			cmd.Writer,
			fx.Invoke(func(
				cfg database.Config,
				logger *zap.SugaredLogger,
			) error {
				gormDB, err := gorm.Open(
					gormpostgres.Open(cfg.CreateDSN()),
					&gorm.Config{
						Logger: gormlogger.New(
							logger,
							gormlogger.Config{
								LogLevel:      gormliblogger.Info,
								SlowThreshold: time.Second * 30,
							},
						),
						DisableAutomaticPing: true,
					},
				)
				if err != nil {
					return err
				}
				g := gen.BuildGenerator(gormDB)
				g.Execute()
				return nil
			}),
		)

		fxApp.Run()

		return fxApp.Err()
	},
}
