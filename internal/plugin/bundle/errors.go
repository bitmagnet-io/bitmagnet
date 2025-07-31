package bundle

import "errors"

var (
	Err                  = errors.New("bundles")
	ErrInvalidRef        = errors.New("invalid ref")
	ErrAlreadyRegistered = errors.New("already registered")
)
