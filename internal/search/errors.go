package search

import "errors"

var (
	Err                   = errors.New("search")
	ErrUnknownIndex       = errors.New("unknown index")
	ErrDuplicateAdapter   = errors.New("duplicate adapter key")
	ErrInvalidAdapterType = errors.New("invalid adapter type")
)
