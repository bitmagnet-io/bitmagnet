package testutils

import (
	"context"
	"errors"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

var (
	RunnerSimple = TestRunner{
		Name:          "simple",
		StartupSleep:  time.Millisecond * 100,
		RunSleep:      time.Minute,
		ShutdownSleep: time.Millisecond * 100,
	}
	RunnerNilShutdowner = TestRunner{
		Name:          "nil shutdowner",
		StartupSleep:  time.Millisecond * 100,
		RunSleep:      time.Minute,
		ShutdownSleep: time.Millisecond * 100,
		ShutdownNil:   true,
	}
	RunnerStartupFailure = TestRunner{
		Name:       "startup failure",
		StartupErr: errors.New("startup failure"),
	}
)

type TestRunner struct {
	Name          string
	StartupSleep  time.Duration
	StartupErr    error
	ShutdownSleep time.Duration
	ShutdownErr   error
	RunSleep      time.Duration
	CancelCause   error
	ShutdownNil   bool
}

func (w *TestRunner) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		select {
		case <-ctx.Done():
		case <-time.After(w.StartupSleep):
		}

		go func() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(w.RunSleep):
			}

			if w.CancelCause != nil {
				cancel(w.CancelCause)
			}
		}()

		var shutdowner runner.Shutdowner

		if !w.ShutdownNil {
			shutdowner = func(ctx context.Context) error {
				select {
				case <-ctx.Done():
				case <-time.After(w.ShutdownSleep):
				}

				return w.ShutdownErr
			}
		}

		return shutdowner, w.StartupErr
	}
}
