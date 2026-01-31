package batch

import "errors"

var (
	Err              = errors.New("batch")
	ErrUninitialized = errors.New("uninitialized")
	ErrFlush         = errors.New("flush failed")
	ErrShutdown      = errors.New("shutdown failed")
)
