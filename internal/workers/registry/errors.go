package registry

import "errors"

var (
	Err                   = errors.New("worker registry")
	ErrUnknownWorker      = errors.New("unknown worker")
	ErrStart              = errors.New("start failed")
	ErrShutdown           = errors.New("shutdown failed")
	ErrRestart            = errors.New("restart failed")
	ErrPartial            = errors.New("partial failure")
	ErrCircularDependency = errors.New("circular dependency detected")
	ErrDependency         = errors.New("dependency failed")
)
