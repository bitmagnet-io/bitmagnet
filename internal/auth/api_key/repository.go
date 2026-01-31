package api_key

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/auth/rbac"
	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"gorm.io/gorm"
)

type Repository interface {
	Create(
		ctx context.Context,
		userID int,
		name string,
		hash []byte,
		permissions []rbac.ObjectAction,
		expiresAt time.Time,
	) (id int, err error)
	Get(ctx context.Context, id int) (model.APIKey, error)
	List(ctx context.Context, req ListRequest) (ListResult, error)
	Delete(ctx context.Context, req DeleteRequest) error
}

func NewRepository(dao database.DaoTransactionProvider) Repository {
	return repository{
		dao: dao,
	}
}

type repository struct {
	dao database.DaoTransactionProvider
}

func (r repository) Create(
	ctx context.Context,
	userID int,
	name string,
	hash []byte,
	permissions []rbac.ObjectAction,
	expiresAt time.Time,
) (id int, err error) {
	dao, err := r.dao.Dao()
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrRepository, err)
	}

	record := model.APIKey{
		UserID: userID,
		Name:   name,
		Hash:   hash,
		Permissions: slice.Map(permissions, func(perm rbac.ObjectAction) model.APIKeyPermission {
			return model.APIKeyPermission{
				Namespace: perm.Namespace,
				Object:    perm.Object,
				Action:    perm.Action,
			}
		}),
	}

	if !expiresAt.IsZero() {
		record.ExpiresAt = sql.NullTime{
			Time:  expiresAt,
			Valid: true,
		}
	}

	err = dao.APIKey.WithContext(ctx).Create(&record)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrRepository, err)
	}

	return record.ID, nil
}

func (r repository) Get(ctx context.Context, id int) (model.APIKey, error) {
	dao, err := r.dao.Dao()
	if err != nil {
		return model.APIKey{}, fmt.Errorf("%w: %w", ErrRepository, err)
	}

	apiKey, err := dao.APIKey.
		WithContext(ctx).
		Where(
			dao.APIKey.ID.Eq(id),
		).
		Preload(
			dao.APIKey.User,
			dao.APIKey.User.Permissions,
			dao.APIKey.Permissions,
		).
		First()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrNotFound
	}

	if err != nil {
		return model.APIKey{}, fmt.Errorf("%w: %w", ErrRepository, err)
	}

	return *apiKey, nil
}

func (r repository) List(ctx context.Context, req ListRequest) (ListResult, error) {
	dao, err := r.dao.Dao()
	if err != nil {
		return ListResult{}, fmt.Errorf("%w: %w", ErrRepository, err)
	}

	limit := 0
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
		return ListResult{}, fmt.Errorf("%w: %w", ErrRepository, err)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	apiKeys, err := q.
		Offset(offset).
		Preload(
			dao.APIKey.User,
			dao.APIKey.User.Permissions,
			dao.APIKey.Permissions,
		).
		Find()
	if err != nil {
		return ListResult{}, fmt.Errorf("%w: %w", ErrRepository, err)
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

func (r repository) Delete(ctx context.Context, req DeleteRequest) error {
	dao, err := r.dao.Dao()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRepository, err)
	}

	info, err := dao.APIKey.WithContext(ctx).
		Where(
			dao.APIKey.UserID.Eq(req.UserID),
			dao.APIKey.ID.Eq(req.APIKeyID),
		).Delete()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRepository, err)
	}

	if info.RowsAffected < 1 {
		return fmt.Errorf("%w: %w", ErrRepository, ErrNotFound)
	}

	return nil
}
