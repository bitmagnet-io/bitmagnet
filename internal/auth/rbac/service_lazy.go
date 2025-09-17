package rbac

import (
	"errors"
	"sync"
)

type ServiceLazy interface {
	Service
	SetService(service Service) error
}

func NewServiceLazy() ServiceLazy {
	return &serviceLazy{}
}

type serviceLazy struct {
	mtx sync.RWMutex
	Service
}

func (e *serviceLazy) SetService(service Service) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()

	if e.Service != nil {
		return errors.New("service already set")
	}

	e.Service = service

	return nil
}
