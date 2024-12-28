package ports

import "styl-monolith/internal/users/core/domain"

type UserPort interface {
	CreateUser(user domain.User) (domain.User, error)
	ListUsers(page, size int) ([]domain.User, int, error)
}

type RolePort interface {
	CreateRole(role domain.Role) (domain.Role, error)
	ListRoles(page, size int) ([]domain.Role, int, error)
}
