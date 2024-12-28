package postgres

import (
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
