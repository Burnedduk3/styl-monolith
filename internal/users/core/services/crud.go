package services

import (
	"styl-monolith/internal/users/core/domain"
)

type CrudServiceStruct struct {
}

type CrudService interface {
	CreateUser(name, email, country, lastIp) domain.User
	CreateRole() domain.Role
	DeleteUser()
	DeleteRole()
	UpdateUser()
	UpdateRole()
	GetUser()
	GetRole()
	ListUsers()
	ListRoles()
}

func NewCrudService() CrudService {
	return &CrudServiceStruct{}
}
