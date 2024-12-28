package models

import (
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/validator"
)

type UserPayload struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
}

func (c *UserPayload) Validate() error {
	if validator.IsStringEmpty(c.Name) ||
		validator.IsStringEmpty(c.Email) ||
		validator.IsStringEmpty(c.Country) ||
		validator.IsStringEmpty(c.CountryCode) ||
		validator.IsStringEmpty(c.Phone) {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestEmptyFields,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestEmptyFields),
			nil)
	}
	if !validator.ValidateStringIsCountryCode(c.CountryCode) ||
		!validator.ValidateStringIsEmail(c.Email) ||
		!validator.ValidateStringIsPhoneNumber(c.Phone) {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestValidationError),
			nil)
	}
	return nil
}

type UserResponse struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	Status      string `json:"status"`
}

func NewCreateUserResponseFromDomainUser(user domain.User) UserResponse {
	return UserResponse{
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		Country:     user.Country,
		Phone:       user.Phone,
		CountryCode: user.CountryCode,
		Status:      string(user.Status),
	}
}

type RolePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *RolePayload) Validate() error {
	if validator.IsStringEmpty(r.Name) ||
		validator.IsStringEmpty(r.Description) {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestEmptyFields,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestEmptyFields),
			nil)
	}
	return nil
}

type RoleResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCreateRoleResponseFromDomainRole(role domain.Role) RoleResponse {
	return RoleResponse{
		Name:        role.Name,
		Description: role.Description,
	}
}
