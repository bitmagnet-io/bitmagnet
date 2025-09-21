package api_key

import (
	"context"
	"fmt"
)

type DeleteRequest struct {
	UserID   int
	APIKeyID int
}

func (s service) Delete(ctx context.Context, req DeleteRequest) error {
	err := s.repository.Delete(ctx, req)
	if err != nil {
		return fmt.Errorf("%w: %w: %w", Err, ErrDelete, err)
	}

	return nil
}
