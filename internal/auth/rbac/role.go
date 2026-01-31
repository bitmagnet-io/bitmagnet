package rbac

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleEditor Role = "editor"
	RoleUser   Role = "user"
	RoleAnon   Role = "anon"
)

func (r Role) String() string {
	return string(r)
}

func CoreRoles() []Role {
	return []Role{
		RoleAdmin,
		RoleEditor,
		RoleUser,
		RoleAnon,
	}
}

type RoleInfo struct {
	Role
	Core        bool
	Permissions []Permission
}
