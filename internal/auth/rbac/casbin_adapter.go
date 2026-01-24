package rbac

import (
	"errors"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

type casbinAdapter struct {
	permissions []Permission
}

func (a *casbinAdapter) setPermissions(perms []Permission) {
	a.permissions = perms
}

func (a *casbinAdapter) LoadPolicy(model model.Model) error {
	for _, perm := range a.permissions {
		tokens := policyTokensFromPermission(perm)
		if err := persist.LoadPolicyArray(tokens, model); err != nil {
			return err
		}
	}

	return nil
}

func policyTokensFromPermission(permission Permission) []string {
	return []string{
		"p",
		subjectString(permission),
		objectString(permission.ObjectAction()),
		permission.ObjectAction().Action,
	}
}

var errUnimplemented = errors.New("unimplemented")

func (casbinAdapter) UpdatePolicy(_ string, _ string, _, _ []string) error {
	return errUnimplemented
}

func (casbinAdapter) UpdatePolicies(_ string, _ string, _, _ [][]string) error {
	return errUnimplemented
}

func (casbinAdapter) UpdateFilteredPolicies(
	_ string,
	_ string,
	_ [][]string,
	_ int,
	_ ...string,
) ([][]string, error) {
	return nil, errUnimplemented
}

func (casbinAdapter) SavePolicy(model.Model) error {
	return errUnimplemented
}

func (casbinAdapter) AddPolicy(string, string, []string) error {
	return errUnimplemented
}

func (casbinAdapter) AddPolicies(string, string, [][]string) error {
	return errUnimplemented
}

func (casbinAdapter) RemovePolicy(string, string, []string) error {
	return errUnimplemented
}

func (casbinAdapter) RemovePolicies(string, string, [][]string) error {
	return errUnimplemented
}

func (casbinAdapter) RemoveFilteredPolicy(string, string, int, ...string) error {
	return errUnimplemented
}
