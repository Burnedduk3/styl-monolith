package postgres

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"styl-monolith/internal/users/adapters/postgres/models"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/ports"
)

type UserRepository struct {
	log  *logrus.Logger
	conn *gorm.DB
}

func NewUserRepository(log *logrus.Logger, conn *gorm.DB) ports.UserPort {
	return &UserRepository{log: log, conn: conn}
}

func (r *UserRepository) CreateUser(user domain.User) (domain.User, error) {
	rUser := models.NewPostgresUserFromDomainUser(user)
	result := r.conn.Create(&rUser)
	if result.Error != nil {
		fmt.Println(result.Error)
		return domain.User{}, fmt.Errorf(
			"error creating user: %s",
			result.Error,
		)
	}
	return rUser.ToUserDomain(), nil
}
