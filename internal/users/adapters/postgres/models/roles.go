package models

import (
	"gorm.io/gorm"
	"styl-monolith/internal/users/core/domain"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
}

func NewPostgresRoleFromDomainRole(role domain.Role) Role {
	return Role{
		Name:        role.Name,
		Description: role.Description,
	}
}

func (r *Role) ToRoleDomain() domain.Role {
	mappedRole := domain.Role{
		ID:          r.Model.ID,
		Name:        r.Name,
		Description: r.Description,
		Created:     r.CreatedAt,
		Updated:     r.UpdatedAt,
	}
	return mappedRole
}
