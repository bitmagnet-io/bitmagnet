package api_key

import (
	"context"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

const secretLength = 32

type ListRequest struct {
	UserID int
	Page   int
	Limit  int
	Offset int
}

type ListResult struct {
	APIKeys    []model.APIKey
	TotalCount int
}

type Service interface {
	Create(ctx context.Context, req CreateRequest) (CreateResult, error)
	Auth(ctx context.Context, key string) (model.APIKey, error)
	List(ctx context.Context, req ListRequest) (ListResult, error)
	Delete(ctx context.Context, id int) error
}

func New(dao database.DaoTransactionProvider) Service {
	return &service{
		dao: dao,
	}
}

type service struct {
	dao database.DaoTransactionProvider
}
