package persister

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Adder interface {
	Add(ctx context.Context, payload Input) error
}

type Persister interface {
	runner.Provider
	Adder
}
