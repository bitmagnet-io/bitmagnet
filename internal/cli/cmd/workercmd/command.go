package workercmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/app"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Command = &cli.Command{
	Name:  "worker",
	Usage: "Run or list workers",
	Commands: []*cli.Command{
		{
			Name:  "run",
			Usage: "Runs workers",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "all",
					Value: false,
				},
				&cli.StringSliceFlag{
					Name: "keys",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				signalChan := make(chan os.Signal, 1)
				signal.Notify(signalChan, os.Interrupt, os.Kill)

				appCtx, appCancel := context.WithCancel(ctx)
				defer appCancel()

				var err error

				fxApp := app.New(
					ctx,
					cmd.Writer,
					fx.Invoke(func(
						registry *registry.Registry,
						logger *zap.SugaredLogger,
						fxShutdowner fx.Shutdowner,
					) {
						logger = logger.Named("cli")

						var workerNames []string

						defer func() {
							_ = fxShutdowner.Shutdown()
						}()

						if cmd.Bool("all") {
							workerNames = registry.Workers()
						} else {
							workerNames = cmd.StringSlice("keys")

							if len(workerNames) == 0 {
								err = errors.New("no workers specified")

								return
							}
						}

						started := make(chan error, 1)

						go func() {
							started <- registry.Start(appCtx, workerNames...)
						}()

						select {
						case sig := <-signalChan:
							appCancel()
							err = fmt.Errorf("startup aborted due to signal: %s", sig)
							return
						case err = <-started:
							if err != nil {
								return
							}
						}

						select {
						case sig := <-signalChan:
							switch sig {
							case os.Interrupt:
								logger.Info(
									"received signal interrupt, starting graceful shutdown",
								)
							case os.Kill:
								logger.Warn("received signal kill, exiting immediately")
								cancel()
								return
							}
						}

						shutdownCtx, shutdownCancel := context.WithTimeout(ctx, time.Minute)
						defer shutdownCancel()

						shutdownCtx, shutdownCancel = signal.NotifyContext(
							shutdownCtx,
							os.Interrupt,
							os.Kill,
						)
						defer shutdownCancel()

						err = registry.ShutdownAll(shutdownCtx)

						if err == nil {
							logger.Info("shutdown complete")
						} else {
							logger.Errorw("shutdown failed", "err", err)
						}
					}),
				)

				go fxApp.Run()

				select {
				case <-appCtx.Done():
				case <-fxApp.Done():
				}

				return errors.Join(err, fxApp.Err(), fxApp.Stop(ctx))
			},
		},
		{
			Name:  "list",
			Usage: "Lists available workers",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fxApp := app.New(
					ctx,
					cmd.Writer,
					fx.Invoke(func(
						registry *registry.Registry,
					) error {
						for _, w := range registry.Workers() {
							if _, err := fmt.Fprintln(cmd.Writer, w); err != nil {
								return err
							}
						}

						return nil
					}),
				)

				fxApp.Run()

				return fxApp.Err()
			},
		},
	},
}
