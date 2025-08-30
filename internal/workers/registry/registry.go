package registry

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/ref"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/zap"
)

type Registry struct {
	workers ref.Map[*worker.Worker]
	logger  *zap.Logger
}

func NewRegistry(logger *zap.Logger, options ...Option) (*Registry, error) {
	r := &Registry{
		workers: ref.NewMap[*worker.Worker](),
		logger:  logger,
	}

	for _, option := range options {
		option(r)
	}

	return r, r.validate()
}

func (r *Registry) validate() error {
	var errs []error

	for _, ref := range r.workers.Refs() {
		for _, dep := range r.dependenciesForward(ref).Refs() {
			if dep.Equals(ref) {
				errs = append(errs, fmt.Errorf("%w: %s", ErrCircularDependency, ref))
			} else if !r.workers.Has(dep) {
				errs = append(errs, fmt.Errorf("%w: %s", ErrUnknownWorker, dep))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %w", Err, errors.Join(errs...))
	}

	return nil
}

func (r *Registry) Runner() runner.Runner {
	return func(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error) {
		err := r.Autostart(ctx)

		return func(ctx context.Context) error {
			return r.ShutdownAll(ctx)
		}, err
	}
}

func (r *Registry) Workers() []ref.Ref {
	return r.workers.Refs()
}

func (r *Registry) WorkersState() ref.Map[WorkerState] {
	result := ref.NewMap[WorkerState]()

	for _, wrk := range r.workers.Values() {
		result.Set(wrk.Ref(), WorkerState{
			StateInfo:  wrk.State(),
			DependsOn:  r.dependenciesForward(wrk.Ref()),
			RequiredBy: r.dependenciesBackward(wrk.Ref()),
			Autostart:  wrk.Autostart(),
		})
	}

	return result
}

func (r *Registry) Start(ctx context.Context, workers ...ref.Ref) error {
	workerMap, err := r.workerMap(true, workers...)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	var (
		mtx            sync.RWMutex
		partialSuccess bool
	)

	workersDone := ref.NewMap[chan struct{}]()
	for _, ref := range workerMap.Refs() {
		workersDone.Set(ref, make(chan struct{}))
	}

	workerErrs := ref.NewMap[error]()

	workerDone := func(ref ref.Ref, err error) {
		mtx.Lock()
		defer mtx.Unlock()

		if err != nil {
			workerErrs.Set(ref, fmt.Errorf("%s: %w", ref, err))
		} else {
			partialSuccess = true
		}

		close(workersDone.Get(ref))
	}

	waitForWorker := func(ref ref.Ref) error {
		ch := workersDone.Get(ref)
		select {
		case <-ch:
			mtx.RLock()
			defer mtx.RUnlock()

			return workerErrs.Get(ref)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	for _, wrk := range workerMap.Values() {
		go func(wrk *worker.Worker) {
			var err error
			defer func() {
				workerDone(wrk.Ref(), err)
			}()

			for _, dep := range wrk.Dependencies() {
				err = waitForWorker(dep)
				if err != nil {
					err = fmt.Errorf("%w: %w", ErrDependency, err)
					return
				}
			}

			_, err = wrk.Start(ctx)
		}(wrk)
	}

	err = errors.Join(slice.Map(workers, waitForWorker)...)
	if err != nil {
		sentinel := ErrStart
		if partialSuccess {
			sentinel = fmt.Errorf("%w: %w", sentinel, ErrPartial)
		}

		return fmt.Errorf("%w: %w: %w", Err, sentinel, err)
	}

	return err
}

func (r *Registry) Autostart(ctx context.Context) error {
	return r.Start(ctx, slice.Filter(r.workers.Refs(), func(ref ref.Ref) bool {
		return r.workers.Get(ref).Autostart()
	})...)
}

func (r *Registry) workerMap(forward bool, workers ...ref.Ref) (ref.Map[*worker.Worker], error) {
	result := ref.NewMap[*worker.Worker]()

	for _, wRef := range workers {
		wrk, ok := r.workers.GetOK(wRef)
		if !ok {
			return result, fmt.Errorf("%w: %s", ErrUnknownWorker, wRef)
		}

		result.Set(wRef, wrk)

		deps := ref.NewSet()
		if forward {
			deps = r.dependenciesForward(wRef)
		} else {
			deps = r.dependenciesBackward(wRef)
		}

		for _, ref := range deps.Refs() {
			dep, ok := r.workers.GetOK(ref)
			if !ok {
				return result, fmt.Errorf("%w: %s", ErrUnknownWorker, ref)
			}

			result.Set(ref, dep)
		}
	}

	return result, nil
}

func (r *Registry) Shutdown(ctx context.Context, workers ...ref.Ref) error {
	workerMap, err := r.workerMap(false, workers...)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	var (
		mtx            sync.RWMutex
		partialSuccess bool
	)

	workersDone := ref.NewMap[chan struct{}]()
	for _, ref := range workerMap.Refs() {
		workersDone.Set(ref, make(chan struct{}))
	}

	workerErrs := ref.NewMap[error]()

	workerDone := func(ref ref.Ref, err error) {
		mtx.Lock()
		defer mtx.Unlock()

		if err != nil {
			workerErrs.Set(ref, fmt.Errorf("%s: %w", ref, err))
		} else {
			partialSuccess = true
		}

		close(workersDone.Get(ref))
	}

	waitForWorker := func(ref ref.Ref) error {
		ch := workersDone.Get(ref)
		select {
		case <-ch:
			mtx.RLock()
			defer mtx.RUnlock()

			return workerErrs.Get(ref)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	for _, wrk := range workerMap.Values() {
		go func(wrk *worker.Worker) {
			var errs []error
			defer func() {
				workerDone(wrk.Ref(), errors.Join(errs...))
			}()

			for _, dep := range r.dependenciesBackward(wrk.Ref()).Refs() {
				err = waitForWorker(dep)
				if err != nil {
					errs = append(errs, fmt.Errorf("%w: %w", ErrDependency, err))
				}
			}

			err = wrk.Shutdown(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}(wrk)
	}

	err = errors.Join(slice.Map(workers, waitForWorker)...)
	if err != nil {
		sentinel := ErrShutdown
		if partialSuccess {
			sentinel = fmt.Errorf("%w: %w", sentinel, ErrPartial)
		}

		err = fmt.Errorf("%w: %w: %w", Err, sentinel, err)
	}

	return err
}

func (r *Registry) ShutdownAll(ctx context.Context) error {
	return r.Shutdown(ctx, r.Workers()...)
}

func (r *Registry) Restart(ctx context.Context, workers ...ref.Ref) error {
	workerMap, err := r.workerMap(false, workers...)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	refs := workerMap.Refs()

	shutdownErr := r.Shutdown(ctx, refs...)
	startErr := r.Start(ctx, refs...)

	err = errors.Join(shutdownErr, startErr)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrRestart, err)
	}

	return nil
}

func (r *Registry) RestartAll(ctx context.Context) error {
	var workersToRestart []ref.Ref

	for _, wrk := range r.workers.Values() {
		if wrk.State().State != worker.StateIdle {
			workersToRestart = append(workersToRestart, wrk.Ref())
		}
	}

	return r.Restart(ctx, workersToRestart...)
}

func (r *Registry) dependenciesForward(rf ref.Ref) ref.Set {
	result := ref.NewSet()

	wrk, ok := r.workers.GetOK(rf)
	if ok {
		for _, dep := range wrk.Dependencies() {
			result.Set(dep, struct{}{})
			for _, childDep := range r.dependenciesForward(dep).Refs() {
				result.Set(childDep, struct{}{})
			}
		}
	}

	return result
}

func (r *Registry) dependenciesBackward(rf ref.Ref) ref.Set {
	result := ref.NewSet()

	for _, wrk := range r.workers.Values() {
		if wrk.DependsOn(rf) {
			result.Set(wrk.Ref(), struct{}{})
			for _, childDep := range r.dependenciesBackward(wrk.Ref()).Refs() {
				result.Set(childDep, struct{}{})
			}
		}
	}

	return result
}
