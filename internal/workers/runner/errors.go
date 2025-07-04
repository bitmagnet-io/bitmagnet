package runner

import "errors"

var (
	Err                  = errors.New(Namespace)
	ErrCompleted         = errors.New("completed")
	ErrShutdownRequested = errors.New("shutdown requested")
	ErrAlreadyRunning    = errors.New("already running")
)
