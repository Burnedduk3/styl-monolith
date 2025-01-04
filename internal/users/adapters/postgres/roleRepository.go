package postgres

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"styl-monolith/internal/users/adapters/postgres/models"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/pkg/errorhandler"
)

type RoleRepository struct {
	log  *logrus.Logger
	conn *gorm.DB
}

func NewRoleRepository(log *logrus.Logger, conn *gorm.DB) *RoleRepository {
	return &RoleRepository{log: log, conn: conn}
}

// CreateRole adds a new role to the database and returns the created domain.Role or an error if the operation fails.
func (r *RoleRepository) CreateRole(role domain.Role) (domain.Role, error) {
	rRole := models.NewPostgresRoleFromDomainRole(role)
	result := r.conn.Create(&rRole)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrRoleDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Role{}, domainError
	}
	return rRole.ToRoleDomain(), nil
}

// ListRoles retrieves a paginated list of roles from the database and calculates the total number of pages for the results.
func (r *RoleRepository) ListRoles(page, size int) ([]domain.Role, int, error) {
	var roles []models.Role
	var dRoles []domain.Role
	offset := (page - 1) * size
	var totalCount int64
	r.conn.Model(&models.Role{}).Count(&totalCount)
	result := r.conn.Limit(size).Offset(offset).Find(&roles)
	// Calculate total pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return nil, 0, domainError
	}
	for _, role := range roles {
		dRoles = append(dRoles, role.ToRoleDomain())
	}
	return dRoles, totalPages, nil
}

// GetRoleByID retrieves a role by its ID from the database and maps it to a domain.Role or returns an error.
func (r *RoleRepository) GetRoleByID(id uint) (domain.Role, error) {
	var role models.Role
	result := r.conn.First(&role, id)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrRoleDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Role{}, domainError
	}
	return role.ToRoleDomain(), nil
}

// UpdateRole updates an existing role in the database and returns the updated role or an error if the operation fails.
func (r *RoleRepository) UpdateRole(role domain.Role) (domain.Role, error) {
	rRole := models.NewPostgresRoleFromDomainRole(role)
	result := r.conn.Model(&models.Role{}).Where("id = ?", role.ID).Save(&rRole)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrRoleDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Role{}, domainError
	}
	return rRole.ToRoleDomain(), nil
}

// DeleteRole deletes a role by its ID from the database and returns an error if the operation fails or the role is not found.
func (r *RoleRepository) DeleteRole(id uint) error {
	result := r.conn.Delete(&models.Role{}, id)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrRoleDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleDatabaseUnableToCompleteOperation),
			result.Error)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrRoleNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrRoleNotFound), id),
			nil,
		)
	}
	return nil
}

// GetRoleByName retrieves a role using its name from the database and returns a domain.Role or an error if not found.
func (r *RoleRepository) GetRoleByName(name string) (domain.Role, error) {
	var role models.Role
	result := r.conn.First(&role, "name = ?", name)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Role{}, errorhandler.NewDomainError(
				errorhandler.ErrRoleNotFound,
				fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrRoleNotFound), name),
				nil,
			)
		}
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrRoleDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Role{}, domainError
	}
	return role.ToRoleDomain(), nil
}

// PartialRoleUpdate updates specified fields of a role identified by its ID and returns the updated role or an error.
func (r *RoleRepository) PartialRoleUpdate(id uint, updatedRole domain.Role) (domain.Role, error) {
	rRole := models.NewPostgresRoleFromDomainRole(updatedRole)
	result := r.conn.Model(&models.Role{}).Where("id = ?", id).Updates(rRole)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrRoleDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.Role{}, domainError
	}
	return rRole.ToRoleDomain(), nil
}
