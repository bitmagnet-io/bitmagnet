package runner

import (
	"context"
	"errors"
	"fmt"
)

type Provider interface {
	Runner() Runner
}

type Runner func(ctx context.Context, cancel context.CancelCauseFunc) (Shutdowner, error)

func (r Runner) Runner() Runner {
	return r
}

func (r Runner) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	shutdowner, err := r(ctx, cancel)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, errors.Join(err, shutdowner(ctx)))
	}

	<-ctx.Done()

	cause := context.Cause(ctx)
	if !errors.Is(cause, ErrCompleted) && !errors.Is(cause, ErrShutdownRequested) {
		return cause
	}

	return nil
}

func SimpleRunner(fn func(context.Context) error) Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (Shutdowner, error) {
		go func() {
			err := fn(ctx)
			if err == nil {
				err = ErrCompleted
			}

			cancel(err)
		}()

		return func(context.Context) error {
			cancel(ErrShutdownRequested)

			return nil
		}, nil
	}
}

type Shutdowner func(context.Context) error

func (s Shutdowner) Call(ctx context.Context) error {
	if s == nil {
		return nil
	}

	return s(ctx)
}
