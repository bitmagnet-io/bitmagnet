package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/cmd"
	"github.com/bitmagnet-io/bitmagnet/internal/plugin/app"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/pkg/env"
	"github.com/bitmagnet-io/bitmagnet/pkg/fs"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Cmd     = cmd.Cmd
	Command = cmd.Command

	App[Deps any] struct {
		builder app.Builder
		logger  *zap.Logger
	}

	pluginOps[Deps any] interface {
		NewRunner(RunnerProvider[Deps]) Runner
		buildPlugin(func(*App[Deps]))
		log(lvl zapcore.Level, msg string, fields ...zap.Field)
		wrapSignals(
			sigProvider env.SignalsProvider,
			rnr runner.Runner,
			shutdownTimeout time.Duration,
		) runner.Runner
	}

	PluginCommand[Deps any] interface {
		cmd.Command
		pluginOps[Deps]
	}

	Runner func(env.Env) error

	RunnerProvider[Deps any] func(Deps) runner.Runner
)

func NewApp[Deps any](appBuilder app.Builder) App[Deps] {
	return App[Deps]{
		builder: appBuilder,
	}
}

func (cmd *App[Deps]) NewRunner(runnerProvider RunnerProvider[Deps]) Runner {
	return func(envEnv env.Env) error {
		fxApp := cmd.builder.Build(
			envOption(envEnv),
			fx.Invoke(func(deps Deps, logger *zap.Logger, shutdowner fx.Shutdowner) error {
				cmd.buildPlugin(func(cmd *App[Deps]) {
					cmd.logger = logger.Named("app")
				})

				err := cmd.wrapSignals(
					envEnv,
					runnerProvider(deps),
					time.Second*30,
				).Run(envEnv)

				return errors.Join(err, shutdowner.Shutdown())
			}),
		)

		if err := fxApp.Err(); err != nil {
			return err
		}

		go fxApp.Run()

		<-fxApp.Done()

		return fxApp.Stop(envEnv)
	}
}

func (cmd *App[Deps]) buildPlugin(fn func(*App[Deps])) {
	fn(cmd)
}

func (cmd *App[Deps]) log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	cmd.logger.Log(lvl, msg, fields...)
}

func (cmd *App[Deps]) wrapSignals(
	sigProvider env.SignalsProvider,
	rnr runner.Runner,
	shutdownTimeout time.Duration,
) runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		var (
			shutdownOnce sync.Once
			shutdowner   runner.Shutdowner
			err          error
		)

		started := make(chan error, 1)

		go func() {
			shutdowner, err = rnr(ctx, cancel)

			started <- err
		}()

		sigChan := sigProvider.Signals(os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-sigChan:
			err = fmt.Errorf("startup aborted due to signal: %s", sig)
		case err = <-started:
		}

		if err != nil {
			cancel(err)

			return runner.NopShutdowner, err
		}

		doShutdown := func(ctx context.Context) error {
			var err error

			shutdownOnce.Do(func() {
				err = shutdowner(ctx)

				cancel(runner.ErrShutdownRequested)
			})

			return err
		}

		go func() {
			select {
			case <-ctx.Done():
				return
			case sig := <-sigChan:
				switch sig {
				case os.Interrupt:
					cmd.log(zapcore.InfoLevel, "received signal, starting graceful shutdown", zap.String("signal", sig.String()))
				case syscall.SIGTERM:
					cmd.log(zapcore.WarnLevel, "received signal, exiting immediately", zap.String("signal", sig.String()))
					cancel(fmt.Errorf("received signal: %s", sig))

					return
				}
			}

			shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
			defer shutdownCancel()

			doShutdown(shutdownCtx)
		}()

		return func(shutdownCtx context.Context) error {
			return doShutdown(shutdownCtx)
		}, nil
	}
}

func envOption(envEnv env.Env) fx.Option {
	return fx.Supply(
		fx.Annotate(
			envEnv,
			fx.As(new(env.Env)),
			fx.As(new(env.Context)),
			fx.As(new(env.Stdin)),
			fx.As(new(env.Stdout)),
			fx.As(new(env.Stderr)),
			fx.As(new(env.VarsLookup)),
			fx.As(new(env.ArgsProvider)),
			fx.As(new(env.SignalsProvider)),
			fx.As(new(fs.FSProvider)),
		),
	)
}
