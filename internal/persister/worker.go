package persister

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"go.uber.org/zap"
)

type worker struct {
	*flusher
	mtx      sync.RWMutex
	in       chan Input
	shutdown chan struct{}
	maxSize  int
	maxWait  time.Duration
	logger   *zap.Logger
}

func (w *worker) Add(ctx context.Context, payload Input) error {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	in := w.in

	if in == nil {
		return fmt.Errorf("%w: %w", Err, runner.ErrNotRunning)
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", Err, ctx.Err())
	case <-w.shutdown:
		return fmt.Errorf("%w: %w", Err, runner.ErrShutdownRequested)
	case in <- payload:
		return nil
	}
}

func (w *worker) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		w.mtx.Lock()
		defer w.mtx.Unlock()

		if w.in != nil {
			return runner.NopShutdowner, fmt.Errorf("%w: %w", Err, runner.ErrAlreadyRunning)
		}

		in := make(chan Input, 1)

		w.in = in

		shutdown := make(chan struct{})

		payload := newPayload()

		lastFlush := time.Now()

		checkFlush := func(ctx context.Context) error {
			size := payload.len()
			if size == 0 {
				lastFlush = time.Now()
				return nil
			}

			if !payload.shouldFlush &&
				size < w.maxSize &&
				time.Since(lastFlush) < w.maxWait {
				return nil
			}

			job := persistJob{
				flusher: w.flusher,
				payload: *payload,
				stats:   make(AllTablesStats),
			}

			payload = newPayload()

			w.logger.Debug("flushing", zap.Int("size", size))

			err := job.run(ctx)

			if err != nil {
				w.logger.Error("flush failed", zap.Error(err))

				return err
			}

			w.logger.Info("flushed", job.stats.LogFields()...)

			lastFlush = time.Now()

			return nil
		}

		go func() {
			defer cancel(nil)

			for {
				select {
				case <-ctx.Done():
					return
				case <-shutdown:
					return
				case <-time.After(w.maxWait - time.Since(lastFlush)):
					_ = checkFlush(ctx)
				case item, ok := <-w.in:
					if !ok {
						return
					}

					item(payload)
					_ = checkFlush(ctx)
				}
			}
		}()

		return func(shutdownCtx context.Context) error {
			close(shutdown)

			w.mtx.Lock()
			w.in = nil
			close(in)
			w.mtx.Unlock()

			select {
			case <-ctx.Done():
				InputFlush(payload)

				return checkFlush(shutdownCtx)
			case <-shutdownCtx.Done():
				return shutdownCtx.Err()
			}
		}, nil
	}
}
