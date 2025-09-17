package api_key

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

func (s *service) List(ctx context.Context, req ListRequest) (ListResult, error) {
	dao, err := s.dao.Dao()
	if err != nil {
		return ListResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := 0

	if req.Offset > 0 {
		offset = req.Offset
	}

	if req.Page > 0 {
		offset += limit * (req.Page - 1)
	}

	q := dao.APIKey.WithContext(ctx)

	if req.UserID > 0 {
		q = q.Where(dao.APIKey.UserID.Eq(req.UserID))
	}

	totalCount, err := q.Count()
	if err != nil {
		return ListResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	apiKeys, err := q.Limit(limit).
		Offset(offset).
		Preload(
			dao.APIKey.User,
			dao.APIKey.Permissions,
		).
		Find()
	if err != nil {
		return ListResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	return ListResult{
		APIKeys: slice.Map(apiKeys, func(apiKey *model.APIKey) model.APIKey {
			apiKey.Hash = nil
			apiKey.User.Password = nil

			return *apiKey
		}),
		TotalCount: int(totalCount),
	}, nil
}
