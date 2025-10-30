package models

// Role представляет роли пользователей.
type Role string

const (
	RoleAdmin   Role = "admin"
	RoleUser    Role = "user"
	RoleManager Role = "manager"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleManager:
		return true
	}
	return false
}
