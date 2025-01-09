package domain

import (
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/validator"
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDeleted  Status = "deleted"
)

type User struct {
	ID          uint
	Username    string
	Password    string
	Name        string
	LastName    string
	Birthday    string
	Email       string
	Country     string
	Phone       string
	CountryCode string
	Role        Role
	Status      Status
	LastIp      string
	LastLogin   time.Time
	Created     time.Time
	Updated     time.Time
	Deleted     time.Time
}

func (u User) UserToUserPayload() UserPayload {
	return UserPayload{
		Name:        u.Name,
		Password:    u.Password,
		LastName:    u.LastName,
		Username:    u.Username,
		Birthday:    u.Birthday,
		Email:       u.Email,
		Country:     u.Country,
		CountryCode: u.CountryCode,
		Phone:       u.Phone,
	}
}

type UserPayload struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Birthday    string `json:"birthday"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
}

func (c *UserPayload) ToUserDomain() User {
	return User{
		Name:        c.Name,
		Password:    c.Password,
		LastName:    c.LastName,
		Username:    c.Username,
		Birthday:    c.Birthday,
		Email:       c.Email,
		Country:     c.Country,
		CountryCode: c.CountryCode,
		Phone:       c.Phone,
	}
}

func (c *UserPayload) Validate() error {
	if validator.IsStringEmpty(c.Name) ||
		validator.IsStringEmpty(c.LastName) ||
		validator.IsStringEmpty(c.Birthday) ||
		validator.IsStringEmpty(c.Password) ||
		validator.IsStringEmpty(c.Email) ||
		validator.IsStringEmpty(c.Country) ||
		validator.IsStringEmpty(c.CountryCode) ||
		validator.IsStringEmpty(c.Phone) {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestEmptyFields,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestEmptyFields),
			nil)
	}
	isValidContryCode := validator.ValidateStringIsCountryCode(c.CountryCode)
	isValidPhone := validator.ValidateStringIsPhoneNumber(c.Phone)
	isValidEmail := validator.ValidateStringIsEmail(c.Email)
	if !isValidContryCode ||
		!isValidPhone ||
		!isValidEmail {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestValidationError),
			nil)
	}
	return nil
}
