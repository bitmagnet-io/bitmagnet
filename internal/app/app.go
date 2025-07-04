package app

import (
	"context"
	"io"

	"github.com/bitmagnet-io/bitmagnet/internal"
	"github.com/bitmagnet-io/bitmagnet/internal/app/appfx"
	"github.com/bitmagnet-io/bitmagnet/internal/logging/loggingfx"
	"go.uber.org/fx"
)

func New(
	ctx context.Context,
	stdout io.Writer,
	opts ...fx.Option,
) *fx.App {
	return fx.New(
		append([]fx.Option{
			appfx.New(),
			loggingfx.WithLogger(),
			fx.Supply(
				fx.Annotate(ctx, fx.As(new(internal.BackgroundContext))),
				fx.Annotate(stdout, fx.As(new(internal.Stdout))),
			),
		}, opts...)...,
	)
}
