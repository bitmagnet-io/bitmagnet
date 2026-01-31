package indexer

import (
	"context"
	"errors"
)

type Multi []Indexer

func (i Multi) Add(ctx context.Context, input Input) error {
	errChan := make(chan error)

	for _, indexer := range i {
		go func(indexer Indexer) {
			errChan <- indexer.Add(ctx, input)
		}(indexer)
	}

	errs := make([]error, 0, len(i))

	for err := range errChan {
		errs = append(errs, err)

		if len(errs) == len(i) {
			break
		}
	}

	return errors.Join(errs...)
}
