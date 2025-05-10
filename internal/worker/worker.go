package worker

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"sync"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RegistryParams struct {
	fx.In
	fx.Shutdowner
	Workers    []Worker    `group:"workers"`
	Decorators []Decorator `group:"worker_decorators"`
	Logger     *zap.SugaredLogger
}

type RegistryResult struct {
	fx.Out
	Registry Registry
}

func NewRegistry(p RegistryParams) (RegistryResult, error) {
	r := &registry{
		mutex:   &sync.RWMutex{},
		workers: make(map[string]Worker),
		logger:  p.Logger,
	}
	for _, w := range p.Workers {
		r.workers[w.Key()] = w
	}

	for _, d := range p.Decorators {
		if err := r.decorate(d.Key, d.Decorate); err != nil {
			return RegistryResult{}, err
		}
	}

	return RegistryResult{Registry: r}, nil
}

type Registry interface {
	Workers() []Worker
	Enable(names ...string) error
	Disable(names ...string) error
	EnableAll()
	DisableAll()
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	decorate(name string, fn DecorateFunction) error
}

type Worker interface {
	Key() string
	Enabled() bool
	Started() bool
	_hook() fx.Hook
	setEnabled(enabled bool)
	setStarted(started bool)
	decorate(DecorateFunction) Worker
}

type DecorateFunction func(fx.Hook) fx.Hook

type Decorator struct {
	Key      string
	Decorate DecorateFunction
}

type worker struct {
	key     string
	hook    fx.Hook
	enabled bool
	started bool
}

func NewWorker(key string, hook fx.Hook) Worker {
	return &worker{
		key:  key,
		hook: hook,
	}
}

func (w *worker) Key() string {
	return w.key
}

func (w *worker) Enabled() bool {
	return w.enabled
}

func (w *worker) Started() bool {
	return w.started
}

func (w *worker) decorate(fn DecorateFunction) Worker {
	return &worker{
		key: w.key,
		hook: fn(fx.Hook{
			OnStart: w.hook.OnStart,
			OnStop:  w.hook.OnStop,
		}),
	}
}

func (w *worker) _hook() fx.Hook {
	return w.hook
}

func (w *worker) setEnabled(enabled bool) {
	w.enabled = enabled
}

func (w *worker) setStarted(started bool) {
	w.started = started
}

type registry struct {
	mutex   *sync.RWMutex
	workers map[string]Worker
	logger  *zap.SugaredLogger
}

func (r *registry) Workers() []Worker {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	keys := slices.Collect(maps.Keys(r.workers))

	sort.Strings(keys)

	return slice.Map(keys, func(s string) Worker {
		return r.workers[s]
	})
}

func (r *registry) Enable(names ...string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, name := range names {
		w, ok := r.workers[name]
		if !ok {
			return fmt.Errorf("worker %s not found", name)
		}

		w.setEnabled(true)
	}

	return nil
}

func (r *registry) Disable(names ...string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, name := range names {
		w, ok := r.workers[name]
		if !ok {
			return fmt.Errorf("worker %s not found", name)
		}

		w.setEnabled(false)
	}

	return nil
}

func (r *registry) EnableAll() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, w := range r.workers {
		w.setEnabled(true)
	}
}

func (r *registry) DisableAll() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, w := range r.workers {
		w.setEnabled(false)
	}
}

var ErrNoWorkersEnabled = errors.New("no workers enabled")

func (r *registry) Start(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	i := 0

	for _, w := range r.workers {
		if w.Enabled() {
			if w.Started() {
				return fmt.Errorf("worker %s already started", w.Key())
			}

			startHook := w._hook().OnStart
			if startHook != nil {
				if err := startHook(ctx); err != nil {
					r.logger.Errorw("error starting worker", "key", w.Key(), "error", err)
					return err
				}
			}

			w.setStarted(true)
			r.logger.Infow("started worker", "key", w.Key())

			i++
		}
	}

	if i == 0 {
		return ErrNoWorkersEnabled
	}

	return nil
}

func (r *registry) Stop(ctx context.Context) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, w := range r.workers {
		if w.Started() {
			stopHook := w._hook().OnStop
			if stopHook != nil {
				if err := stopHook(ctx); err != nil {
					r.logger.Errorw("error stopping worker", "key", w.Key(), "error", err)
					continue
				}
			}

			w.setStarted(false)
			r.logger.Infow("stopped worker", "key", w.Key())
		}
	}

	return nil
}

func (r *registry) decorate(name string, fn DecorateFunction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if w, ok := r.workers[name]; ok {
		r.workers[name] = w.decorate(fn)
		return nil
	}

	return fmt.Errorf("worker %s not found", name)
}
