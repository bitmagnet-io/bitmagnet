package rbac

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/slice"
	"github.com/casbin/casbin/v2"
)

type Enforcer interface {
	Enforce(ctx context.Context, subject Subject, objectAction ObjectAction) (bool, error)
	EnforceAny(ctx context.Context, subjects []Subject, objectAction ObjectAction) (bool, error)
}

type Service interface {
	Enforcer
	Repository
	GetObjectActions() []ObjectAction
}

func NewService(
	repository Repository,
	objectActionProvider ObjectActionProvider,
	permissionProvider PermissionProvider,
	ttl CacheTTL,
) Service {
	return &service{
		sem:                  make(chan struct{}, 1),
		repository:           repository,
		objectActionProvider: objectActionProvider,
		permissionProvider:   permissionProvider,
		ttl:                  time.Duration(ttl),
	}
}

// service implements the Service interface backed by casbin.
// The design is partly working around the fact that casbin does not support context.
type service struct {
	sem                  chan struct{}
	ttl                  time.Duration
	lastUpdate           time.Time
	repository           Repository
	permissionProvider   PermissionProvider
	objectActionProvider ObjectActionProvider
	*casbinDeps
}

type casbinDeps struct {
	*casbin.Enforcer
	*casbinAdapter
}

func (s *service) Enforce(ctx context.Context, subject Subject, objectAction ObjectAction) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case s.sem <- struct{}{}:
	}

	defer func() { <-s.sem }()

	casbin, err := s.acquireCasbin(ctx)
	if err != nil {
		return false, err
	}

	return casbin.Enforce(subjectString(subject), objectString(objectAction), objectAction.Action)
}

func (s *service) EnforceAny(ctx context.Context, subjects []Subject, objectAction ObjectAction) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case s.sem <- struct{}{}:
	}

	defer func() { <-s.sem }()

	casbin, err := s.acquireCasbin(ctx)
	if err != nil {
		return false, err
	}

	result, err := casbin.BatchEnforce(batchRequests(subjects, objectAction))

	return slices.Contains(result, true), err
}

func (s *service) GetAllRoles(ctx context.Context) ([]RoleInfo, error) {
	return s.getRoles(ctx, nil)
}

func (s *service) GetRole(ctx context.Context, role Role) (RoleInfo, error) {
	roles, err := s.getRoles(ctx, []Role{role})
	if err != nil {
		return RoleInfo{}, err
	}

	return roles[0], nil
}

func (s *service) GetRoles(ctx context.Context, roles []Role) ([]RoleInfo, error) {
	if len(roles) == 0 {
		return nil, nil
	}

	return s.getRoles(ctx, roles)
}

func (s *service) getRoles(ctx context.Context, roles []Role) ([]RoleInfo, error) {
	var (
		roleInfos []RoleInfo
		err       error
	)

	if roles == nil {
		roleInfos, err = s.repository.GetAllRoles(ctx)
	} else {
		roleInfos, err = s.repository.GetRoles(ctx, roles)
	}

	if err != nil {
		return nil, err
	}

	seenRoles := make(map[Role]struct{}, len(roleInfos))
	for _, role := range roleInfos {
		seenRoles[role.Role] = struct{}{}
	}

	for _, coreRole := range CoreRoles() {
		if roles != nil && !slices.Contains(roles, coreRole) {
			continue
		}

		if _, ok := seenRoles[coreRole]; !ok {
			roleInfos = append(roleInfos, RoleInfo{
				Role: coreRole,
				Core: true,
			})

			seenRoles[coreRole] = struct{}{}
		}
	}

	if roles != nil {
		missingRoles := slice.FlatMap(roles, func(role Role) []string {
			_, ok := seenRoles[role]

			if !ok {
				return []string{string(role)}
			}

			return nil
		})

		if len(missingRoles) > 0 {
			return nil, fmt.Errorf("roles not found: %s", strings.Join(missingRoles, ", "))
		}
	}

	roleInfos = slice.Map(roleInfos, s.mergeCoreRolePermissions)

	slices.SortFunc(roleInfos, func(a, b RoleInfo) int {
		return cmp.Compare(a.Role, b.Role)
	})

	return roleInfos, nil
}

func (s *service) GetPermissions(ctx context.Context) ([]Permission, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case s.sem <- struct{}{}:
	}

	defer func() { <-s.sem }()

	casin, err := s.acquireCasbin(ctx)
	if err != nil {
		return nil, err
	}

	return casin.permissions, nil
}

func (s *service) PutRole(ctx context.Context, role Role, objectActions []ObjectAction) (RoleInfo, error) {
	select {
	case <-ctx.Done():
		return RoleInfo{}, ctx.Err()
	case s.sem <- struct{}{}:
	}

	defer func() { <-s.sem }()

	roleInfo, err := s.repository.PutRole(ctx, role, objectActions)
	if err != nil {
		return RoleInfo{}, err
	}

	if !s.lastUpdate.IsZero() {
		err = s.updatePermissions(ctx)
		if err != nil {
			return RoleInfo{}, err
		}
	}

	return s.mergeCoreRolePermissions(roleInfo), nil
}

func (s *service) DeleteRole(ctx context.Context, role Role) error {
	if slices.Contains(CoreRoles(), role) {
		return errors.New("core roles cannot be deleted")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.sem <- struct{}{}:
	}

	defer func() { <-s.sem }()

	err := s.repository.DeleteRole(ctx, role)

	if !s.lastUpdate.IsZero() {
		err = s.updatePermissions(ctx)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *service) GetObjectActions() []ObjectAction {
	return s.objectActionProvider()
}

// acquireCasbin returns the casbin enforcer instance, initializing it if necessary.
func (s *service) acquireCasbin(ctx context.Context) (*casbinDeps, error) {
	if s.casbinDeps == nil {
		permissions, err := s.repository.GetPermissions(ctx)
		if err != nil {
			return nil, err
		}

		adapterDynamic := &casbinAdapter{
			permissions: append(s.permissionProvider(), permissions...),
		}

		eCasbin, err := newCasbinEnforcer(adapterDynamic)
		if err != nil {
			return nil, err
		}

		s.casbinDeps = &casbinDeps{
			casbinAdapter: adapterDynamic,
			Enforcer:      eCasbin,
		}

		s.lastUpdate = time.Now()
	} else if time.Since(s.lastUpdate) > s.ttl {
		if err := s.updatePermissions(ctx); err != nil {
			return nil, err
		}
	}

	return s.casbinDeps, nil
}

func (s *service) updatePermissions(ctx context.Context) error {
	permissions, err := s.getPermissions(ctx)
	if err != nil {
		return err
	}

	s.setPermissions(permissions)

	err = s.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}

	s.lastUpdate = time.Now()

	return nil
}

func (s *service) getPermissions(ctx context.Context) ([]Permission, error) {
	perms := make(map[string]map[ObjectAction]Permission)

	repoPerms, err := s.repository.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}

	applyPerm := func(perm Permission) {
		subj := subjectString(perm)
		if _, ok := perms[subj]; !ok {
			perms[subj] = make(map[ObjectAction]Permission)
		}

		perms[subj][perm.ObjectAction()] = perm
	}

	for _, perm := range repoPerms {
		applyPerm(perm)
	}

	for _, perm := range s.permissionProvider() {
		applyPerm(perm)
	}

	var finalPerms []Permission

	for _, objAct := range perms {
		for _, perm := range objAct {
			finalPerms = append(finalPerms, perm)
		}
	}

	sortPermissions(finalPerms)

	return finalPerms, nil
}

func subjectString(sub Subject) string {
	return fmt.Sprintf("%s::%s", sub.SubjectType(), sub.SubjectName())
}

func objectString(objAct ObjectAction) string {
	return fmt.Sprintf("%s::%s", objAct.Namespace, objAct.Object)
}

func batchRequests(subs []Subject, objAct ObjectAction) [][]any {
	return slice.Map(subs, func(sub Subject) []any {
		return []any{
			subjectString(sub),
			objectString(objAct),
			objAct.Action,
		}
	})
}

func (s *service) mergeCoreRolePermissions(roleInfo RoleInfo) RoleInfo {
	permMap := make(map[ObjectAction]bool)
	for _, perm := range roleInfo.Permissions {
		permMap[perm.ObjectAction()] = false
	}

	for _, perm := range s.permissionProvider() {
		if perm.SubjectType() == SubjectTypeRole && perm.SubjectName() == string(roleInfo.Role) {
			permMap[perm.ObjectAction()] = true
		}
	}

	perms := make([]Permission, 0, len(permMap))
	for objAct, core := range permMap {
		perms = append(perms, permission{
			Subject: SubjectRole{
				Role: roleInfo.Role,
			},
			objectAction: objAct,
			core:         core,
		})
	}

	sortPermissions(perms)

	return RoleInfo{
		Role:        roleInfo.Role,
		Core:        slices.Contains(CoreRoles(), roleInfo.Role),
		Permissions: perms,
	}
}

func sortPermissions(perms []Permission) {
	slices.SortFunc(perms, func(a, b Permission) int {
		r := cmp.Compare(subjectString(a), subjectString(b))
		if r != 0 {
			return r
		}

		return a.ObjectAction().Compare(b.ObjectAction())
	})
}
