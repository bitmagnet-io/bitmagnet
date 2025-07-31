package concat

import "errors"

var (
	Err                  = errors.New("concat")
	ErrPartial           = errors.New("partial failure")
	ErrShutdown          = errors.New("shutdown failed")
	ErrAllRunnersStopped = errors.New("all runners stopped")
)
