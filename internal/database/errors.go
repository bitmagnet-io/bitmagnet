package database

import "errors"

var (
	Err              = errors.New(Namespace)
	ErrUninitialized = errors.New("uninitialized")
	ErrPingFailed    = errors.New("ping failed")
)
