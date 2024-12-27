package models

import (
	"gorm.io/gorm"
	"styl-monolith/internal/users/core/domain"
)

type Role struct {
	gorm.Model
	Name string
}

func NewPostgresRoleFromDomainRole(user domain.Role) Role {
	return Role{
		Name: user.Name,
	}
}

func (r *Role) toRoleDomain() domain.Role {
	mappedUser := domain.Role{
		ID:   r.Model.ID,
		Name: r.Name,
	}
	return mappedUser
}
