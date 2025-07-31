package channel

import "errors"

var (
	Err              = errors.New("channel")
	ErrUninitialized = errors.New("uninitialized")
	ErrItem          = errors.New("item failed")
	ErrShutdown      = errors.New("shutdown failed")
)
