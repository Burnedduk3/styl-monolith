package services

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/ports"
	"styl-monolith/pkg/errorhandler"
)

type CrudServiceStruct struct {
	logger       *logrus.Logger
	userCrudRepo ports.UserPort
	roleCrudRepo ports.RolePort
}

type CrudService interface {
	CreateUser(domain.User) (domain.User, error)
	CreateRole(role domain.Role) (domain.Role, error)
	DeleteUserById(id uint) error
	DeleteRoleById(id uint) error
	UpdateUserById(id uint, user domain.User) (domain.User, error)
	UpdateRoleById(id uint, role domain.Role) (domain.Role, error)
	PartialUpdateUserById(id uint, user domain.User) (domain.User, error)
	PartialUpdateRoleById(id uint, role domain.Role) (domain.Role, error)
	GetUserById(id uint) (domain.User, error)
	GetRoleById(id uint) (domain.Role, error)
	GetRoleByName(name string) (domain.Role, error)
	GetUserByUsername(username string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetUserByPhone(phone string) (domain.User, error)
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
	err := c.roleCrudRepo.DeleteRole(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CrudServiceStruct) DeleteUserById(id uint) error {
	err := c.userCrudRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *CrudServiceStruct) GetRoleById(id uint) (domain.Role, error) {
	role, err := c.roleCrudRepo.GetRoleByID(id)
	if err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

func (c *CrudServiceStruct) GetRoleByName(name string) (domain.Role, error) {
	user, err := c.roleCrudRepo.GetRoleByName(name)
	if err != nil {
		return domain.Role{}, err
	}
	return user, nil
}

func (c *CrudServiceStruct) GetUserById(id uint) (domain.User, error) {
	user, err := c.userCrudRepo.GetUserByID(id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (c *CrudServiceStruct) GetUserByEmail(email string) (domain.User, error) {
	user, err := c.userCrudRepo.GetUserByEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (c *CrudServiceStruct) GetUserByUsername(username string) (domain.User, error) {
	user, err := c.userCrudRepo.GetUserByUsername(username)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (c *CrudServiceStruct) GetUserByPhone(phone string) (domain.User, error) {
	user, err := c.userCrudRepo.GetUserByUsername(phone)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
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

func (c *CrudServiceStruct) UpdateRoleById(id uint, role domain.Role) (domain.Role, error) {
	oldRole, err := c.GetRoleById(id)
	if err != nil {
		return domain.Role{}, err
	}
	if role.Name == "" || role.Description == "" {

		return domain.Role{}, errorhandler.NewDomainError(
			errorhandler.ErrRoleIncompleteParamsForCompleteUpdate,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleIncompleteParamsForCompleteUpdate),
			nil,
		)
	}
	oldRole.Name = role.Name
	oldRole.Description = role.Description
	newRole, err := c.roleCrudRepo.UpdateRole(oldRole)
	if err != nil {
		return domain.Role{}, err
	}
	return newRole, nil
}

func (c *CrudServiceStruct) UpdateUserById(id uint, user domain.User) (domain.User, error) {
	oldUser, err := c.GetUserById(id)
	if err != nil {
		return domain.User{}, err
	}
	if oldUser.Name == "" || oldUser.Phone == "" || oldUser.Username == "" || oldUser.Email == "" || oldUser.Country == "" || oldUser.CountryCode == "" {
		return domain.User{}, errorhandler.NewDomainError(
			errorhandler.ErrUserIncompleteParamsForCompleteUpdate,
			errorhandler.GetErrorMessage(errorhandler.ErrUserIncompleteParamsForCompleteUpdate),
			nil,
		)
	}
	oldUser.Name = user.Name
	oldUser.Phone = user.Phone
	oldUser.Username = user.Username
	oldUser.Email = user.Email
	oldUser.Country = user.Country
	oldUser.CountryCode = user.CountryCode
	newRole, err := c.userCrudRepo.UpdateUser(oldUser)
	if err != nil {
		return domain.User{}, err
	}
	return newRole, nil
}

func (c *CrudServiceStruct) PartialUpdateRoleById(id uint, role domain.Role) (domain.Role, error) {
	oldRole, err := c.GetRoleById(id)
	if err != nil {
		return domain.Role{}, err
	}
	if role.Name != "" {
		oldRole.Name = role.Name
	}
	if role.Description != "" {
		oldRole.Description = role.Description
	}
	newRole, err := c.roleCrudRepo.PartialRoleUpdate(id, oldRole)
	if err != nil {
		return domain.Role{}, err
	}
	return newRole, nil
}

func (c *CrudServiceStruct) PartialUpdateUserById(id uint, user domain.User) (domain.User, error) {
	oldUser, err := c.GetUserById(id)

	if err != nil {
		return domain.User{}, err
	}

	if user.Name != "" {
		oldUser.Name = user.Name
	}
	if user.Phone != "" {
		oldUser.Phone = user.Phone
	}
	if user.Username != "" {
		oldUser.Username = user.Username
	}
	if user.Email != "" {
		oldUser.Email = user.Email
	}
	if user.Country != "" {
		oldUser.Country = user.Country
	}
	if user.CountryCode != "" {
		oldUser.CountryCode = user.CountryCode
	}

	newRole, err := c.userCrudRepo.PartialUserUpdate(id, oldUser)
	if err != nil {
		return domain.User{}, err
	}
	return newRole, nil
}
