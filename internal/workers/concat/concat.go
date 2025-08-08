package concat

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

// Runners combines zero or more Runner functions into a single Runner.
// Each Runner is invoked in the provided order, with Shutdowner functions invoked in reverse order.
func Runners(providers ...runner.Provider) runner.Runner {
	return runner.Runner(func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		var (
			shutdowners []runner.Shutdowner
			startupErrs []error
			cancelErrs  []error
			mtx         sync.Mutex
			waitGroup   sync.WaitGroup
		)

		waitGroup.Add(len(providers))

		for _, provider := range providers {
			runCtx, childCancel := context.WithCancelCause(ctx)

			var once sync.Once

			runCancel := func(err error) {
				once.Do(func() {
					childCancel(err)

					if err != nil {
						mtx.Lock()
						cancelErrs = append(cancelErrs, err)
						mtx.Unlock()
					}

					waitGroup.Done()
				})
			}

			shutdowner, err := provider.Runner()(runCtx, runCancel)
			if err != nil {
				runCancel(err)

				startupErrs = append(startupErrs, err)
			}

			shutdowners = append(shutdowners, runner.OnceShutdowner(func(ctx context.Context) error {
				err := shutdowner.Call(ctx)

				// Cancel the worker's run context after graceful shutdown:
				runCancel(nil)

				return err
			}))
		}

		go func() {
			// Cancel the parent context once all worker run contexts have been canceled:
			waitGroup.Wait()

			err := ErrAllRunnersStopped
			cancelErr := errors.Join(cancelErrs...)

			if cancelErr != nil {
				err = fmt.Errorf("%w: %w", err, cancelErr)
			}

			err = fmt.Errorf("%w: %w", Err, err)

			cancel(err)
		}()

		startupErr := errors.Join(startupErrs...)

		if startupErr != nil {
			if len(startupErrs) < len(providers) {
				startupErr = fmt.Errorf("%w: %w", ErrPartial, startupErr)
			}

			startupErr = fmt.Errorf("%w: %w", Err, startupErr)
		}

		return runner.OnceShutdowner(func(ctx context.Context) error {
			var shutdownErrs []error

			for i := len(shutdowners) - 1; i >= 0; i-- {
				shutdownErrs = append(shutdownErrs, shutdowners[i].Call(ctx))
			}

			shutdownErr := errors.Join(shutdownErrs...)
			if shutdownErr != nil {
				if len(shutdownErrs) < len(providers) {
					shutdownErr = fmt.Errorf("%w: %w", ErrPartial, shutdownErr)
				}

				shutdownErr = fmt.Errorf("%w: %w: %w", Err, ErrShutdown, shutdownErr)
			}

			return shutdownErr
		}), startupErr
	})
}
