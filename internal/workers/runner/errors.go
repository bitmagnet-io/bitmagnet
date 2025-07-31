package runner

import "errors"

var (
	Err                  = errors.New("runner")
	ErrCompleted         = errors.New("completed")
	ErrShutdownRequested = errors.New("shutdown requested")
	ErrAlreadyRunning    = errors.New("already running")
	ErrNotRunning        = errors.New("not running")
)
