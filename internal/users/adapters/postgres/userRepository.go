package postgres

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"styl-monolith/internal/users/adapters/postgres/models"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/ports"
	"styl-monolith/pkg/errorhandler"
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
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.User{}, domainError
	}
	return rUser.ToUserDomain(), nil
}

func (r *UserRepository) ListUsers(page, size int) ([]domain.User, int, error) {
	var users []models.User
	var dUsers []domain.User
	offset := (page - 1) * size
	var totalCount int64
	r.conn.Model(&models.User{}).Count(&totalCount)
	result := r.conn.Limit(size).Offset(offset).Find(&users)
	// Calculate total pages
	totalPages := int((totalCount + int64(size) - 1) / int64(size))
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return nil, 0, domainError
	}
	for _, user := range users {
		dUsers = append(dUsers, user.ToUserDomain())
	}
	return dUsers, totalPages, nil
}
