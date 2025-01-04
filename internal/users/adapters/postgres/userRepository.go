package postgres

import (
	"fmt"
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

// CreateUser creates a new user in the database and returns the created user or an error if the operation fails.
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

// ListUsers retrieves a paginated list of users from the database.
// Returns a slice of domain.User, total pages, and an error if any issue occurs during the operation.
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

// GetUserByID retrieves a user by their unique ID from the database.
// Returns a domain.User and an error if the operation fails.
func (r *UserRepository) GetUserByID(id uint) (domain.User, error) {
	var user models.User
	result := r.conn.First(&user, id)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.User{}, domainError
	}
	return user.ToUserDomain(), nil
}

// UpdateUser updates an existing user in the database and returns the updated user or an error if the operation fails.
func (r *UserRepository) UpdateUser(user domain.User) (domain.User, error) {
	rUser := models.NewPostgresUserFromDomainUser(user)
	result := r.conn.Model(&models.User{}).Where("id = ?", user.ID).Save(&rUser)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.User{}, domainError
	}
	return rUser.ToUserDomain(), nil
}

// DeleteUser removes a user record from the database using the provided user ID.
// Returns an error if the operation fails or if the user is not found.
func (r *UserRepository) DeleteUser(id uint) error {
	result := r.conn.Delete(&models.User{}, id)
	if result.Error != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
	}
	if result.RowsAffected == 0 {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserNotFound,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrUserNotFound), id),
			nil,
		)
	}
	return nil
}

// GetUserByEmail retrieves a user from the database by their email address. Returns the user or an error if not found.
func (r *UserRepository) GetUserByEmail(email string) (domain.User, error) {
	var user models.User
	result := r.conn.First(&user, "email = ?", email)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.User{}, domainError
	}
	return user.ToUserDomain(), nil
}

// GetUserByPhone retrieves a user record from the database based on the provided phone number. Returns the user or an error.
func (r *UserRepository) GetUserByPhone(phone string) (domain.User, error) {
	var user models.User
	result := r.conn.First(&user, "phone = ?", phone)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.User{}, domainError
	}
	return user.ToUserDomain(), nil
}

// GetUserByUsername retrieves a user from the database using the provided username.
// Returns a domain.User and an error if the operation fails.
func (r *UserRepository) GetUserByUsername(username string) (domain.User, error) {
	var user models.User
	result := r.conn.First(&user, "username = ?", username)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.User{}, domainError
	}
	return user.ToUserDomain(), nil
}

// PartialUserUpdate updates specific fields of an existing user identified by their ID in the database.
// Returns the updated user as a domain.User or an error if the operation fails.
func (r *UserRepository) PartialUserUpdate(id uint, updatedUser domain.User) (domain.User, error) {
	rUser := models.NewPostgresUserFromDomainUser(updatedUser)
	result := r.conn.Model(&models.User{}).Where("id = ?", id).Updates(rUser)
	if result.Error != nil {
		domainError := errorhandler.NewDomainError(
			errorhandler.ErrUserDatabaseUnableToCompleteOperation,
			errorhandler.GetErrorMessage(errorhandler.ErrUserDatabaseUnableToCompleteOperation),
			result.Error)
		return domain.User{}, domainError
	}
	return rUser.ToUserDomain(), nil
}
