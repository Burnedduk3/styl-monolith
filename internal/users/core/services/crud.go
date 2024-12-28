package services

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/ports"
)

type CrudServiceStruct struct {
	logger       *logrus.Logger
	userCrudRepo ports.UserPort
	roleCrudRepo ports.RolePort
}

type CrudService interface {
	CreateUser(name, username, email, country, phone, countryCode string) (domain.User, error)
	CreateRole(name, description string) (domain.Role, error)
	DeleteUser(id uint) error
	DeleteRole(id uint) error
	UpdateUser(id uint, name, email, country, lastIp string) (domain.User, error)
	UpdateRole(id uint, name string) (domain.Role, error)
	GetUser(id uint) (domain.User, error)
	GetRole(id uint) (domain.Role, error)
	ListUsers(offset, limit int) ([]domain.User, int, error)
	ListRoles(offset, limit int) ([]domain.Role, int, error)
}

func NewCrudService(log *logrus.Logger, userCrudRepo ports.UserPort, roleRepo ports.RolePort) CrudService {
	return &CrudServiceStruct{
		logger:       log,
		userCrudRepo: userCrudRepo,
		roleCrudRepo: roleRepo,
	}
}

func (c *CrudServiceStruct) CreateUser(name, username, email, country, phone, countryCode string) (domain.User, error) {
	domainUser := domain.User{
		Name:        name,
		Username:    username,
		Email:       email,
		Country:     country,
		Phone:       phone,
		CountryCode: countryCode,
	}
	dUser, err := c.userCrudRepo.CreateUser(domainUser)
	if err != nil {
		return domain.User{}, err
	}
	return dUser, nil
}

func (c *CrudServiceStruct) CreateRole(name, description string) (domain.Role, error) {
	domainRole := domain.Role{
		Name:        name,
		Description: description,
	}
	dRole, err := c.roleCrudRepo.CreateRole(domainRole)
	if err != nil {
		return domain.Role{}, err
	}
	return dRole, nil
}

func (c *CrudServiceStruct) DeleteUser(id uint) error {
	// TODO: Add logic for deleting a user by ID.
	return nil
}

func (c *CrudServiceStruct) DeleteRole(id uint) error {
	// TODO: Add logic for deleting a role by ID.
	return nil
}

func (c *CrudServiceStruct) UpdateUser(id uint, name, email, country, lastIp string) (domain.User, error) {
	// TODO: Add logic for updating an existing user's details.
	return domain.User{}, nil
}

func (c *CrudServiceStruct) UpdateRole(id uint, name string) (domain.Role, error) {
	// TODO: Add logic for updating an existing role's details.
	return domain.Role{}, nil
}

func (c *CrudServiceStruct) GetUser(id uint) (domain.User, error) {
	// TODO: Add logic for retrieving a user by their ID.
	return domain.User{}, nil
}

func (c *CrudServiceStruct) GetRole(id uint) (domain.Role, error) {
	// TODO: Add logic for retrieving a role by their ID.
	return domain.Role{}, nil
}

func (c *CrudServiceStruct) ListUsers(page, size int) ([]domain.User, int, error) {
	users, totalPages, err := c.userCrudRepo.ListUsers(page, size)
	if err != nil {
		return []domain.User{}, 0, err
	}
	return users, totalPages, nil
}

func (c *CrudServiceStruct) ListRoles(page, size int) ([]domain.Role, int, error) {
	roles, totalPages, err := c.roleCrudRepo.ListRoles(page, size)
	if err != nil {
		return []domain.Role{}, 0, err
	}
	return roles, totalPages, nil
}
