package api_key

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
)

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

func (s service) List(ctx context.Context, req ListRequest) (ListResult, error) {
	result, err := s.repository.List(ctx, req)
	if err != nil {
		return ListResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	return result, nil
}
