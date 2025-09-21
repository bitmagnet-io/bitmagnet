package rbac

import (
	"fmt"
)

type SubjectType string

const (
	SubjectTypeRole       SubjectType = "role"
	SubjectTypePermission SubjectType = "permission"
)

type Subject interface {
	SubjectType() SubjectType
	SubjectName() string
}

type SubjectRole struct {
	Role Role
}

func (SubjectRole) SubjectType() SubjectType {
	return SubjectTypeRole
}

func (s SubjectRole) SubjectName() string {
	return string(s.Role)
}

type SubjectPermission struct {
	ObjectAction ObjectAction
}

func (SubjectPermission) SubjectType() SubjectType {
	return SubjectTypePermission
}

func (s SubjectPermission) SubjectName() string {
	return fmt.Sprintf("%s::%s::%s", s.ObjectAction.Namespace, s.ObjectAction.Object, s.ObjectAction.Action)
}
