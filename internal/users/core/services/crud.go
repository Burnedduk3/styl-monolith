package services

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/domain"
)

type CrudServiceStruct struct {
	logger *logrus.Logger
}

type CrudService interface {
	CreateUser(name, email, country, lastIp string) domain.User
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

func NewCrudService(log *logrus.Logger) CrudService {
	return &CrudServiceStruct{
		logger: log,
	}
}

func (c *CrudServiceStruct) CreateUser(name, email, country, lastIp string) domain.User {
	return domain.User{}
}

func (c *CrudServiceStruct) CreateRole() domain.Role {
	return domain.Role{}
}

func (c *CrudServiceStruct) DeleteUser() {}

func (c *CrudServiceStruct) DeleteRole() {}

func (c *CrudServiceStruct) UpdateUser() {}

func (c *CrudServiceStruct) UpdateRole() {}

func (c *CrudServiceStruct) GetUser() {}

func (c *CrudServiceStruct) GetRole() {}

func (c *CrudServiceStruct) ListUsers() {}

func (c *CrudServiceStruct) ListRoles() {}
