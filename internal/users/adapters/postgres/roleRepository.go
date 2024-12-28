package postgres

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"styl-monolith/internal/users/adapters/postgres/models"
	"styl-monolith/internal/users/core/domain"
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
		fmt.Println(result.Error)
		return domain.Role{}, fmt.Errorf(
			"error creating user: %s",
			result.Error,
		)
	}
	return rRole.ToRoleDomain(), nil
}
