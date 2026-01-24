package gqlmodel

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/target"
)

type TargetMutation struct {
	Targets target.Registry
}

func (t TargetMutation) Send(ctx context.Context, input target.Params) (target.Result, error) {
	return t.Targets.Send(ctx, input)
}
