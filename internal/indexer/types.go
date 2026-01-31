package indexer

import (
	"context"
)

type Indexer interface {
	Add(context.Context, Input) error
}

type Input func(*payload)
