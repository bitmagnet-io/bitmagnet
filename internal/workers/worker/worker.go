package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/zap"
)

type StateInfo struct {
	Ref ref.Ref
	State
	Err error
}

type Worker struct {
	ref        ref.Ref
	mtx        sync.RWMutex
	state      State
	runner     runner.Provider
	nextState  chan struct{}
	shutdowner runner.Shutdowner
	err        error
	logger     *zap.Logger
	dependsOn  ref.Set
	autostart  bool
	shortLived bool
}

func NewWorker(rf ref.Ref, runner runner.Provider, options ...Option) *Worker {
	wrk := &Worker{
		ref:       rf,
		state:     StateIdle,
		runner:    runner,
		dependsOn: ref.NewSet(),
	}

	for _, option := range options {
		option(wrk)
	}

	if wrk.logger == nil {
		wrk.logger = zap.NewNop()
	}

	return wrk
}

func (w *Worker) Ref() ref.Ref {
	return w.ref
}

func (w *Worker) State() StateInfo {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return StateInfo{
		Ref:   w.ref,
		State: w.state,
		Err:   w.err,
	}
}

func (w *Worker) Err() error {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	return w.err
}

func noopCanceler(error) {}

func noopShutdowner(context.Context) error {
	return nil
}

func (w *Worker) Start(ctx context.Context) (runner.Shutdowner, error) {
	w.mtx.Lock()

	switch w.state {
	case StateStartup:
		nextState := w.nextState
		w.mtx.Unlock()

		<-nextState

		return w.Start(ctx)
	case StateRunning:
		shutdowner := w.shutdowner
		w.mtx.Unlock()

		return shutdowner, nil
	case StateShutdown:
		ch := w.nextState
		w.mtx.Unlock()

		<-ch

		return w.Start(ctx)
	case StateIdle, StateError:
		w.err = nil
		w.state = StateStartup
		nextState := make(chan struct{})
		w.nextState = nextState

		w.mtx.Unlock()

		var closeOnce sync.Once

		doClose := func() {
			closeOnce.Do(func() {
				close(nextState)

				w.nextState = nil
			})
		}

		runCtx, runCancel := context.WithCancelCause(ctx)

		if w.shortLived {
			w.logger.Info("starting")
		} else {
			w.logger.Debug("starting")
		}

		shutdown, err := w.runner.Runner()(runCtx, func(err error) {
			isShutdownRequested := errors.Is(err, runner.ErrShutdownRequested)
			isEndedWithError := err != nil && !isShutdownRequested

			if isEndedWithError {
				w.logger.Error("ended with error", zap.Error(err))
			}

			sentinel := fmt.Errorf("%w: %w", Err, ErrStopped)
			if err == nil {
				err = sentinel
			} else {
				err = fmt.Errorf("%w: %w", sentinel, err)
			}

			w.mtx.Lock()

			if isEndedWithError {
				w.state = StateError
				w.err = err
			} else {
				w.state = StateIdle
				w.err = nil
			}

			doClose()
			runCancel(err)
			w.mtx.Unlock()
		})

		if err == nil {
			if w.shortLived {
				w.logger.Debug("started")
			} else {
				w.logger.Info("started")
			}
		} else {
			w.logger.Error("failed to start", zap.Error(err))
		}

		w.mtx.Lock()

		shutdowner := w.newShutdowner(runCancel, shutdown)
		w.shutdowner = shutdowner

		doClose()

		if err != nil {
			err = fmt.Errorf("%w: %w", ErrStart, err)

			w.state = StateError
			w.err = err

			w.mtx.Unlock()

			err = fmt.Errorf("%w: %w", Err, err)

			runCancel(err)

			return nil, err
		}

		switch w.state {
		case StateStartup:
			w.state = StateRunning
		case StateIdle:
			w.logger.Info("completed")
		}

		w.mtx.Unlock()

		return shutdowner, nil
	default:
		state := w.state
		err := w.err
		w.mtx.Unlock()

		return nil, fmt.Errorf("%w: %w: %w: %s: %w", Err, ErrStart, ErrInvalidState, state.String(), err)
	}
}

func (w *Worker) Shutdown(ctx context.Context) error {
	w.mtx.RLock()
	shutdowner := w.shutdowner
	w.mtx.RUnlock()

	if shutdowner == nil {
		shutdowner = w.newShutdowner(noopCanceler, noopShutdowner)
	}

	return shutdowner(ctx)
}

func (w *Worker) newShutdowner(runCancel context.CancelCauseFunc, shutdown runner.Shutdowner) runner.Shutdowner {
	return func(ctx context.Context) error {
		defer runCancel(nil)

		w.mtx.Lock()

		switch w.state {
		case StateIdle:
			w.mtx.Unlock()

			return nil
		case StateStartup, StateShutdown:
			nextState := w.nextState
			w.mtx.Unlock()

			<-nextState
			w.mtx.RLock()
			shutdowner := w.shutdowner
			err := w.err
			w.mtx.RUnlock()

			if shutdowner != nil {
				err = errors.Join(shutdowner(ctx), err)
			}

			if err != nil {
				err = fmt.Errorf("%w: %w: %w", Err, ErrShutdownFailed, err)
			}

			return err
		case StateRunning:
			w.state = StateShutdown
			w.shutdowner = nil
			nextState := make(chan struct{})
			w.nextState = nextState

			w.mtx.Unlock()

			w.logger.Debug("shutting down")

			var err error

			if shutdown != nil {
				shutdownCtx, shutdownCancel := context.WithCancel(ctx)

				err = shutdown(shutdownCtx)

				shutdownCancel()
			}

			if err == nil {
				w.logger.Info("stopped")
			} else {
				w.logger.Error("shutdown failed", zap.Error(err))
			}

			w.mtx.Lock()

			w.nextState = nil

			if err == nil {
				w.state = StateIdle
			} else {
				w.state = StateError
				err = fmt.Errorf("%w: %w", ErrShutdownFailed, err)
				w.err = err
				err = fmt.Errorf("%w: %w", Err, err)
			}

			w.mtx.Unlock()

			close(nextState)

			return err
		case StateError:
			err := w.err
			nextState := make(chan struct{})
			w.nextState = nextState
			w.state = StateIdle
			w.err = nil
			w.mtx.Unlock()

			close(nextState)

			return fmt.Errorf("%w: %w: %w", Err, ErrShutdownFailed, err)
		default:
			state := w.state
			w.mtx.Unlock()

			return fmt.Errorf("%w: %w: %w: %s", Err, ErrShutdownFailed, ErrInvalidState, state.String())
		}
	}
}

func (w *Worker) Dependencies() []ref.Ref {
	return w.dependsOn.Refs()
}

func (w *Worker) DependsOn(ref ref.Ref) bool {
	return w.dependsOn.Has(ref)
}

func (w *Worker) Autostart() bool {
	return w.autostart
}
