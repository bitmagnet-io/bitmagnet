package periodic

import "errors"

var (
	Err         = errors.New("periodic")
	ErrInvoke   = errors.New("invoke failed")
	ErrShutdown = errors.New("shutdown failed")
)
