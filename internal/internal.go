package internal

import (
	"context"
	"io"
)

type (
	BackgroundContext context.Context
	Stdout            io.Writer
)
