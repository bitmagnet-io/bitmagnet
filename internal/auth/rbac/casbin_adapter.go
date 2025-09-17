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

func (casbinAdapter) UpdatePolicy(sec string, ptype string, oldRule, newRule []string) error {
	return errUnimplemented
}

func (casbinAdapter) UpdatePolicies(sec string, ptype string, oldRules, newRules [][]string) error {
	return errUnimplemented
}

func (casbinAdapter) UpdateFilteredPolicies(sec string, ptype string, newRules [][]string, fieldIndex int, fieldValues ...string) ([][]string, error) {
	return nil, errUnimplemented
}

func (casbinAdapter) SavePolicy(model model.Model) error {
	return errUnimplemented
}

func (casbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errUnimplemented
}

func (casbinAdapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	return errUnimplemented
}

func (casbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errUnimplemented
}

func (casbinAdapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return errUnimplemented
}

func (casbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errUnimplemented
}
