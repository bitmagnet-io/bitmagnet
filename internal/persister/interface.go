package persister

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/workers/runner"
)

type Adder interface {
	Add(context.Context, Input) error
}

type Persister interface {
	runner.Provider
	Adder
}
