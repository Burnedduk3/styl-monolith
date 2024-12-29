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
	CreateUser(requestUser domain.User) (domain.User, error)
	CreateRole(requestRole domain.Role) (domain.Role, error)
	DeleteUserById(id uint) error
	DeleteRoleById(id uint) error
	UpdateUserById(id uint, requestUser domain.User) (domain.User, error)
	UpdateRoleById(id uint, requestRole domain.Role) (domain.Role, error)
	GetUserById(id uint) (domain.User, error)
	GetRoleById(id uint) (domain.Role, error)
	GetRoleByName(name string) (domain.Role, error)
	GetUserByUsername(username string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
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

func (c *CrudServiceStruct) CreateRole(requestRole domain.Role) (domain.Role, error) {
	dRole, err := c.roleCrudRepo.CreateRole(requestRole)
	if err != nil {
		return domain.Role{}, err
	}
	return dRole, nil
}

func (c *CrudServiceStruct) CreateUser(requestUser domain.User) (domain.User, error) {
	dUser, err := c.userCrudRepo.CreateUser(requestUser)
	if err != nil {
		return domain.User{}, err
	}
	return dUser, nil
}

func (c *CrudServiceStruct) DeleteRoleById(id uint) error {
	// TODO: Add logic for deleting a role by ID.
	return nil
}

func (c *CrudServiceStruct) DeleteUserById(id uint) error {
	// TODO: Add logic for deleting a user by ID.
	return nil
}

func (c *CrudServiceStruct) GetRoleById(id uint) (domain.Role, error) {
	// TODO: Add logic for retrieving a role by their ID.
	return domain.Role{}, nil
}

func (c *CrudServiceStruct) GetRoleByName(name string) (domain.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CrudServiceStruct) GetUserById(id uint) (domain.User, error) {
	// TODO: Add logic for retrieving a user by their ID.
	return domain.User{}, nil
}

func (c *CrudServiceStruct) GetUserByEmail(email string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CrudServiceStruct) GetUserByUsername(username string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CrudServiceStruct) ListRoles(page, size int) ([]domain.Role, int, error) {
	roles, totalPages, err := c.roleCrudRepo.ListRoles(page, size)
	if err != nil {
		return []domain.Role{}, 0, err
	}
	return roles, totalPages, nil
}

func (c *CrudServiceStruct) ListUsers(page, size int) ([]domain.User, int, error) {
	users, totalPages, err := c.userCrudRepo.ListUsers(page, size)
	if err != nil {
		return []domain.User{}, 0, err
	}
	return users, totalPages, nil
}

func (c *CrudServiceStruct) UpdateRoleById(id uint, requestRole domain.Role) (domain.Role, error) {
	// TODO: Add logic for updating an existing role's details.
	return domain.Role{}, nil
}

func (c *CrudServiceStruct) UpdateUserById(id uint, requestUser domain.User) (domain.User, error) {
	// TODO: Add logic for updating an existing user's details.
	return domain.User{}, nil
}
