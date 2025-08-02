package registry

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/worker"
	"go.uber.org/zap"
)

type Registry struct {
	workers map[string]*worker.Worker
	logger  *zap.Logger
}

func NewRegistry(logger *zap.Logger, options ...Option) (*Registry, error) {
	r := &Registry{
		logger:  logger,
		workers: make(map[string]*worker.Worker),
	}

	for _, option := range options {
		option(r)
	}

	return r, r.validate()
}

func (r *Registry) validate() error {
	var errs []error

	for key := range r.workers {
		for dep := range r.dependenciesForward(key) {
			if dep == key {
				errs = append(errs, fmt.Errorf("%w: %s", ErrCircularDependency, key))
			} else if _, ok := r.workers[dep]; !ok {
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

func (r *Registry) Workers() []string {
	workers := slices.Collect(maps.Keys(r.workers))
	slices.Sort(workers)

	return workers
}

func (r *Registry) WorkersState() map[string]WorkerState {
	result := make(map[string]WorkerState, len(r.workers))
	for k, v := range r.workers {
		result[k] = WorkerState{
			StateInfo:  v.State(),
			DependsOn:  r.dependenciesForward(k),
			RequiredBy: r.dependenciesBackward(k),
			Autostart:  r.workers[k].Autostart(),
		}
	}

	return result
}

func (r *Registry) Start(ctx context.Context, workers ...string) error {
	workerMap, err := r.workerMap(true, workers...)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	var (
		mtx            sync.RWMutex
		partialSuccess bool
	)

	workersDone := make(map[string]chan struct{}, len(workerMap))
	for key := range workerMap {
		workersDone[key] = make(chan struct{})
	}

	workerErrs := make(map[string]error)

	workerDone := func(key string, err error) {
		mtx.Lock()
		defer mtx.Unlock()

		if err != nil {
			workerErrs[key] = fmt.Errorf("%s: %w", key, err)
		} else {
			partialSuccess = true
		}

		close(workersDone[key])
	}

	waitForWorker := func(key string) error {
		ch := workersDone[key]
		select {
		case <-ch:
			mtx.RLock()
			defer mtx.RUnlock()

			return workerErrs[key]
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	for key, wrk := range workerMap {
		go func(key string) {
			var err error
			defer func() {
				workerDone(key, err)
			}()

			for _, dep := range wrk.Dependencies() {
				err = waitForWorker(dep)
				if err != nil {
					err = fmt.Errorf("%w: %w", ErrDependency, err)
					return
				}
			}

			_, err = wrk.Start(ctx)
		}(key)
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
	return r.Start(ctx, slice.Filter(r.Workers(), func(key string) bool {
		return r.workers[key].Autostart()
	})...)
}

func (r *Registry) workerMap(forward bool, workers ...string) (map[string]*worker.Worker, error) {
	result := make(map[string]*worker.Worker, len(workers))

	for _, key := range workers {
		wrk, ok := r.workers[key]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownWorker, key)
		}

		result[key] = wrk

		var deps map[string]struct{}
		if forward {
			deps = r.dependenciesForward(key)
		} else {
			deps = r.dependenciesBackward(key)
		}

		for depKey := range deps {
			dep, ok := r.workers[depKey]
			if !ok {
				return nil, fmt.Errorf("%w: %s", ErrUnknownWorker, depKey)
			}

			result[depKey] = dep
		}
	}

	return result, nil
}

func (r *Registry) Shutdown(ctx context.Context, workers ...string) error {
	workerMap, err := r.workerMap(false, workers...)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	var (
		mtx            sync.RWMutex
		partialSuccess bool
	)

	workersDone := make(map[string]chan struct{}, len(workerMap))
	for key := range workerMap {
		workersDone[key] = make(chan struct{})
	}

	workerErrs := make(map[string]error)

	workerDone := func(key string, err error) {
		mtx.Lock()
		defer mtx.Unlock()

		if err != nil {
			workerErrs[key] = fmt.Errorf("%s: %w", key, err)
		} else {
			partialSuccess = true
		}

		close(workersDone[key])
	}

	waitForWorker := func(key string) error {
		ch := workersDone[key]
		select {
		case <-ch:
			mtx.RLock()
			defer mtx.RUnlock()

			return workerErrs[key]
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	for key, wrk := range workerMap {
		go func(key string) {
			var errs []error
			defer func() {
				workerDone(key, errors.Join(errs...))
			}()

			for dep := range r.dependenciesBackward(key) {
				err = waitForWorker(dep)
				if err != nil {
					errs = append(errs, fmt.Errorf("%w: %w", ErrDependency, err))
				}
			}

			err = wrk.Shutdown(ctx)
			if err != nil {
				errs = append(errs, err)
			}
		}(key)
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

func (r *Registry) Restart(ctx context.Context, workers ...string) error {
	workerMap, err := r.workerMap(false, workers...)
	if err != nil {
		return fmt.Errorf("%w: %w", Err, err)
	}

	allWorkers := slices.Collect(maps.Keys(workerMap))

	shutdownErr := r.Shutdown(ctx, allWorkers...)
	startErr := r.Start(ctx, allWorkers...)

	err = errors.Join(shutdownErr, startErr)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrRestart, err)
	}

	return nil
}

func (r *Registry) RestartAll(ctx context.Context) error {
	var workersToRestart []string

	for key, wrk := range r.workers {
		if wrk.State().State != worker.StateIdle {
			workersToRestart = append(workersToRestart, key)
		}
	}

	return r.Restart(ctx, workersToRestart...)
}

func (r *Registry) dependenciesForward(key string) worker.DependencyMap {
	result := make(worker.DependencyMap)

	wrk, ok := r.workers[key]
	if ok {
		for _, dep := range wrk.Dependencies() {
			result[dep] = struct{}{}
			for childDep := range r.dependenciesForward(dep) {
				result[childDep] = struct{}{}
			}
		}
	}

	return result
}

func (r *Registry) dependenciesBackward(key string) worker.DependencyMap {
	result := make(worker.DependencyMap)

	for otherKey, wrk := range r.workers {
		if wrk.DependsOn(key) {
			result[otherKey] = struct{}{}
			for childDep := range r.dependenciesBackward(otherKey) {
				result[childDep] = struct{}{}
			}
		}
	}

	return result
}
