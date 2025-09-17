package user

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type ListUsersParams struct {
	UsernameLike string
	Limit        int
	Page         int
	Offset       int
}

type ListUsersResult struct {
	Users      []model.User
	TotalCount int
}

func (s *service) List(ctx context.Context, params ListUsersParams) (ListUsersResult, error) {
	limit := params.Limit
	if limit < 1 {
		limit = defaultLimit
	}

	page := params.Page
	if page < 1 {
		page = 1
	}

	offset := params.Offset
	offset += (page - 1) * limit

	if offset < 0 {
		offset = 0
	}

	dao, err := s.Dao()
	if err != nil {
		return ListUsersResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	q := dao.User.WithContext(ctx)

	if params.UsernameLike != "" {
		// todo: Check edge cases, e.g. % input
		q = q.Where(dao.User.Username.Like("%" + params.UsernameLike + "%"))
	}

	totalCount, err := q.Count()
	if err != nil {
		return ListUsersResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	orderBy := dao.User.ID.Asc()

	users, err := q.Limit(limit).Offset(offset).Order(orderBy).Find()
	if err != nil {
		return ListUsersResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	return ListUsersResult{
		Users: slice.Map(users, func(user *model.User) model.User {
			return *user
		}),
		TotalCount: int(totalCount),
	}, nil
}
