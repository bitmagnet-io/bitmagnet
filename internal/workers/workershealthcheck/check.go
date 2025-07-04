package workershealthcheck

import (
	"context"
	"errors"
	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
	"sync"
	"time"
)

type WorkersHealthCheck struct {
	mtx              sync.Mutex
	registry         *registry.Registry
	receivedRegistry chan struct{}
}

func New() *WorkersHealthCheck {
	return &WorkersHealthCheck{
		receivedRegistry: make(chan struct{}),
	}
}

func (c *WorkersHealthCheck) Check() health.Check {
	return health.Check{
		Name:    "workers",
		Timeout: time.Second,
		Check: func(ctx context.Context) error {
			<-c.receivedRegistry
			state := c.registry.WorkersState()

			var errs []error
			for _, st := range state {
				errs = append(errs, st.Err)
			}

			return errors.Join(errs...)
		},
	}
}

func (c *WorkersHealthCheck) SetRegistry(r *registry.Registry) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.registry != nil {
		return errors.New("registry already set")
	}

	c.registry = r
	close(c.receivedRegistry)

	return nil
}
