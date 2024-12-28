package ports

import "styl-monolith/internal/users/core/domain"

type UserPort interface {
	CreateUser(user domain.User) (domain.User, error)
}

type RolePort interface {
	CreateRole(role domain.Role) (domain.Role, error)
}
