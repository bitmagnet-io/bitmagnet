package registry

import (
	"errors"
	"sync"
)

type StateProvider interface {
	WorkersState() map[string]WorkerState
}

var _ StateProvider = (*Registry)(nil)

type CircuitBreaker interface {
	StateProvider
	ReceiveRegistry(*Registry) error
}

func NewCircuitBreaker() CircuitBreaker {
	return &circuitBreaker{
		received: make(chan struct{}),
	}
}

type circuitBreaker struct {
	mtx      sync.RWMutex
	received chan struct{}
	registry *Registry
}

func (c *circuitBreaker) WorkersState() map[string]WorkerState {
	<-c.received
	return c.registry.WorkersState()
}

func (c *circuitBreaker) ReceiveRegistry(r *Registry) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.registry != nil {
		return errors.New("registry already set")
	}
	c.registry = r
	close(c.received)
	return nil
}

var _ CircuitBreaker = (*circuitBreaker)(nil)
