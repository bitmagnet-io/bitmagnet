package rbac

import (
	"fmt"
)

type SubjectType string

const (
	SubjectTypeUser SubjectType = "user"
	SubjectTypeRole SubjectType = "role"
)

type Subject interface {
	SubjectType() SubjectType
	SubjectName() string
}

type SubjectUser struct {
	ID int32
}

func (SubjectUser) SubjectType() SubjectType {
	return SubjectTypeUser
}

func (s SubjectUser) SubjectName() string {
	return fmt.Sprintf("%d", s.ID)
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
