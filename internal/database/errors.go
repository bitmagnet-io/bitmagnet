package database

import "errors"

var (
	Err              = errors.New("postgres")
	ErrUninitialized = errors.New("uninitialized")
	ErrPingFailed    = errors.New("ping failed")
)
