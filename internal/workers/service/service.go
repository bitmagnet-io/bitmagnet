package service

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Service interface {
	Name() string
	Run(ctx context.Context, cancel context.CancelCauseFunc) (runner.Shutdowner, error)
}
