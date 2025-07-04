package runner

import (
	"context"
)

func Nop(_ context.Context, _ context.CancelCauseFunc) (Shutdowner, error) {
	return NopShutdowner, nil
}

func NopShutdowner(_ context.Context) error {
	return nil
}
