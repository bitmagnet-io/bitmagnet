package rbac

import (
	"context"
	"slices"

	"github.com/bitmagnet-io/bitmagnet/internal/database"
	"github.com/bitmagnet-io/bitmagnet/internal/database/dao"
	"github.com/bitmagnet-io/bitmagnet/internal/model"
	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"gorm.io/gen"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetAllRoles(ctx context.Context) ([]RoleInfo, error)
	GetRole(ctx context.Context, roles Role) (RoleInfo, error)
	GetRoles(ctx context.Context, roles []Role) ([]RoleInfo, error)
	GetPermissions(ctx context.Context) ([]Permission, error)
	PutRole(ctx context.Context, role Role, objectActions []ObjectAction) (RoleInfo, error)
	DeleteRole(ctx context.Context, role Role) error
}

func NewRepository(dao database.DaoTransactionProvider) Repository {
	return &repository{
		dao: dao,
	}
}

type repository struct {
	dao database.DaoTransactionProvider
}

func (r *repository) GetAllRoles(ctx context.Context) ([]RoleInfo, error) {
	dao, err := r.dao.Dao()
	if err != nil {
		return nil, err
	}

	dbRoles, err := dao.WithContext(ctx).Role.
		Preload(dao.Role.Permissions).
		Find()
	if err != nil {
		return nil, err
	}

	return slice.Map(dbRoles, roleInfoFromModel), nil
}

func (r *repository) GetRole(ctx context.Context, role Role) (RoleInfo, error) {
	dao, err := r.dao.Dao()
	if err != nil {
		return RoleInfo{}, err
	}

	dbRole, err := dao.WithContext(ctx).Role.
		Preload(dao.Role.Permissions).
		Where(dao.Role.Name.Eq(string(role))).
		First()
	if err != nil {
		return RoleInfo{}, err
	}

	return roleInfoFromModel(dbRole), nil
}

func (r *repository) GetRoles(ctx context.Context, roles []Role) ([]RoleInfo, error) {
	dao, err := r.dao.Dao()
	if err != nil {
		return nil, err
	}

	dbRoles, err := dao.WithContext(ctx).Role.
		Preload(dao.Role.Permissions).
		Where(dao.Role.Name.In(slice.Map(roles, func(role Role) string {
			return string(role)
		})...)).
		Find()
	if err != nil {
		return nil, err
	}

	return slice.Map(dbRoles, roleInfoFromModel), nil
}

func (r *repository) GetPermissions(ctx context.Context) ([]Permission, error) {
	dao, err := r.dao.Dao()
	if err != nil {
		return nil, err
	}

	dbPerms, err := dao.WithContext(ctx).RolePermission.Find()
	if err != nil {
		return nil, err
	}

	return slice.Map(dbPerms, permissionFromModel), nil
}

func (r *repository) PutRole(ctx context.Context, role Role, objectActions []ObjectAction) (RoleInfo, error) {
	var roleInfo RoleInfo

	err := r.dao.DaoTransaction(func(tx *dao.Query) error {
		err := tx.WithContext(ctx).
			Role.
			Clauses(clause.OnConflict{
				DoNothing: true,
			}).
			Create(&model.Role{
				Name: string(role),
			})
		if err != nil {
			return err
		}

		if len(objectActions) == 0 {
			return nil
		}

		_, err = tx.WithContext(ctx).
			RolePermission.
			Where(tx.RolePermission.RoleName.Eq(string(role))).
			Delete()
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).
			RolePermission.
			Create(slice.Map(objectActions, func(objAct ObjectAction) *model.RolePermission {
				return &model.RolePermission{
					RoleName:  string(role),
					Namespace: objAct.Namespace,
					Object:    objAct.Object,
					Action:    objAct.Action,
				}
			})...)
		if err != nil {
			return err
		}

		roleModel, err := tx.WithContext(ctx).
			Role.
			Preload(tx.Role.Permissions).
			Where(tx.Role.Name.Eq(string(role))).
			First()
		if err != nil {
			return err
		}

		roleInfo = roleInfoFromModel(roleModel)

		return nil
	})

	return roleInfo, err
}

func (r *repository) DeleteRole(ctx context.Context, role Role) error {
	dao, err := r.dao.Dao()
	if err != nil {
		return err
	}

	_, err = dao.WithContext(ctx).
		Role.
		Where(dao.Role.Name.Eq(string(role))).
		Delete()

	return err
}

func (r *repository) DeleteRolePermissions(
	ctx context.Context,
	role Role,
	objectActions []ObjectAction,
) (RoleInfo, error) {
	var roleInfo RoleInfo

	err := r.dao.DaoTransaction(func(tx *dao.Query) error {
		if len(objectActions) > 0 {
			_, err := tx.WithContext(ctx).RolePermission.Scopes(
				func(scope gen.Dao) gen.Dao {
					return scope.Where(dao.RolePermission.RoleName.Eq(string(role)))
				},
				func(scope gen.Dao) gen.Dao {
					return scope.Or(
						slice.Map(objectActions, func(objAct ObjectAction) gen.Condition {
							return dao.RolePermission.
								Where(dao.RolePermission.Object.Eq(objAct.Namespace)).
								Where(dao.RolePermission.Object.Eq(objAct.Object)).
								Where(dao.RolePermission.Action.Eq(objAct.Action))
						})...)
				},
			).Delete()
			if err != nil {
				return err
			}
		}

		roleModel, err := tx.WithContext(ctx).
			Role.
			Preload(dao.Role.Permissions).
			Where(dao.Role.Name.Eq(string(role))).
			First()
		if err != nil {
			return err
		}

		roleInfo = roleInfoFromModel(roleModel)

		return nil
	})

	return roleInfo, err
}

func permissionFromModel(perm *model.RolePermission) Permission {
	return NewPermission(
		SubjectRole{
			Role: Role(perm.RoleName),
		},
		ObjectAction{
			Namespace: perm.Namespace,
			Object:    perm.Object,
			Action:    perm.Action,
		},
	)
}

func roleInfoFromModel(roleModel *model.Role) RoleInfo {
	role := Role(roleModel.Name)

	return RoleInfo{
		Role: role,
		Core: slices.Contains(CoreRoles(), role),
		Permissions: slice.Map(roleModel.Permissions, func(perm model.RolePermission) Permission {
			return permissionFromModel(&perm)
		}),
	}
}
