package api_key

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (CreateResult, error)
	Auth(ctx context.Context, key string) (model.APIKey, error)
	List(ctx context.Context, req ListRequest) (ListResult, error)
	Delete(ctx context.Context, req DeleteRequest) error
}

func NewService(repository Repository) Service {
	return service{
		repository: repository,
	}
}

type service struct {
	repository Repository
}
