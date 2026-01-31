package workershealthcheck

import (
	"context"
	"errors"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/health"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
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
		Check: func(context.Context) error {
			return errors.Join(
				slice.Map(c.registry.WorkersState().Values(), func(st registry.WorkerState) error {
					return st.Err
				})...)
		},
	}
}
