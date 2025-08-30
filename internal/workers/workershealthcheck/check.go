package workershealthcheck

import (
	"context"
	"errors"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/workers/registry"
)

type WorkersHealthCheck struct {
	registry registry.StateProvider
}

func New(reg registry.StateProvider) *WorkersHealthCheck {
	return &WorkersHealthCheck{
		registry: reg,
	}
}

func (c *WorkersHealthCheck) Check() health.Check {
	return health.Check{
		Name:    "workers",
		Timeout: time.Second,
		Check: func(ctx context.Context) error {
			state := c.registry.WorkersState()

			var errs []error
			for _, st := range state.Values() {
				errs = append(errs, st.Err)
			}

			return errors.Join(errs...)
		},
	}
}
