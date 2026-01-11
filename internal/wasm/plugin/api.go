package plugin

import (
	context "context"
	"errors"

	pool "github.com/jolestar/go-commons-pool/v2"
)

type apiService[T any] struct {
	pool       *pool.ObjectPool
	getService func(module *module) (T, error)
}

func (s *apiService[T]) do(ctx context.Context, f func(T) error) error {
	modAny, err := s.pool.BorrowObject(ctx)
	if err != nil {
		return err
	}

	service, err := s.getService(modAny.(*module))
	if err != nil {
		_ = s.pool.ReturnObject(ctx, modAny)
		return err
	}

	err = f(service)
	retErr := s.pool.ReturnObject(ctx, modAny)

	return errors.Join(err, retErr)
}
