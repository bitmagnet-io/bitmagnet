package worker

import "errors"

var (
	Err               = errors.New(Namespace)
	ErrStart          = errors.New("start failed")
	ErrShutdownFailed = errors.New("shutdown failed")
	ErrStopped        = errors.New("stopped")
)
