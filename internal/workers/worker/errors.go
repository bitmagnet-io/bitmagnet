package worker

import "errors"

var (
	Err               = errors.New("worker")
	ErrStart          = errors.New("start failed")
	ErrShutdownFailed = errors.New("shutdown failed")
	ErrStopped        = errors.New("stopped")
)
