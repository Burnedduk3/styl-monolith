package models

import (
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/validator"
)

type UserPayload struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
}

func (c *UserPayload) ToUserDomain() domain.User {
	return domain.User{
		Name:        c.Name,
		Username:    c.Username,
		Email:       c.Email,
		Country:     c.Country,
		CountryCode: c.CountryCode,
		Phone:       c.Phone,
	}
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
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	Status      string `json:"status"`
}

func NewUserResponseFromDomainUser(user domain.User) UserResponse {
	return UserResponse{
		Id:          user.ID,
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
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *RolePayload) ToRoleDomain() domain.Role {
	return domain.Role{
		Name:        r.Name,
		Description: r.Description,
	}
}

func (r *RolePayload) Validate() error {
	if validator.IsStringEmpty(r.Name) ||
		validator.IsStringEmpty(r.Description) {
		return errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestEmptyFields,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestEmptyFields),
			nil)
	}
	return nil
}

type RoleResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewRoleResponseFromDomainRole(role domain.Role) RoleResponse {
	return RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

type PaginationResponse struct {
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	TotalPages  int         `json:"total_pages"`
	Data        interface{} `json:"data"`
}
