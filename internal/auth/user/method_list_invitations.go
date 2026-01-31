package user

import (
	"context"
	"fmt"

	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
)

type ListInvitationsParams struct {
	Limit  int
	Page   int
	Offset int
}

type ListInvitationsResult struct {
	Invitations []model.Invitation
	TotalCount  int
}

func (s *service) ListInvitations(ctx context.Context, params ListInvitationsParams) (ListInvitationsResult, error) {
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
		return ListInvitationsResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	q := dao.Invitation.WithContext(ctx)

	totalCount, err := q.Count()
	if err != nil {
		return ListInvitationsResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	orderBy := dao.Invitation.CreatedAt.Desc()

	invitations, err := q.
		Preload(
			dao.Invitation.CreatedByUser,
			dao.Invitation.ClaimedByUser,
		).
		Limit(limit).
		Offset(offset).
		Order(orderBy).
		Find()
	if err != nil {
		return ListInvitationsResult{}, fmt.Errorf("%w: %w: %w", Err, ErrList, err)
	}

	return ListInvitationsResult{
		Invitations: slice.Map(invitations, func(invitation *model.Invitation) model.Invitation {
			return *invitation
		}),
		TotalCount: int(totalCount),
	}, nil
}
