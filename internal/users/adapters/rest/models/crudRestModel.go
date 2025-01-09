package models

import (
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/validator"
)

type UserResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Birthday    string `json:"birthday"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	Status      string `json:"status"`
	Role        RoleResponse
}

func NewUserResponseFromDomainUser(user domain.User) UserResponse {
	response := UserResponse{
		Id:          user.ID,
		Name:        user.Name,
		LastName:    user.LastName,
		Birthday:    user.Birthday,
		Username:    user.Username,
		Email:       user.Email,
		Country:     user.Country,
		Phone:       user.Phone,
		CountryCode: user.CountryCode,
		Status:      string(user.Status),
	}
	if user.Role.ID != 0 {
		response.Role = NewRoleResponseFromDomainRole(user.Role)
	}
	return response
}

func (ur *UserResponse) ToDomainUser() domain.User {
	return domain.User{
		ID:          ur.Id,
		Name:        ur.Name,
		LastName:    ur.LastName,
		Birthday:    ur.Birthday,
		Username:    ur.Username,
		Email:       ur.Email,
		Country:     ur.Country,
		CountryCode: ur.CountryCode,
		Phone:       ur.Phone,
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

type UpdateUserRolePayload struct {
	RoleId uint `json:"role_id"`
}
