package models

import (
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/validator"
	"time"
)

type UserResponse struct {
	Id          uint             `json:"id"`
	Name        string           `json:"name"`
	LastName    string           `json:"last_name"`
	Username    string           `json:"username"`
	Birthday    string           `json:"birthday"`
	Email       string           `json:"email"`
	Country     string           `json:"country"`
	Phone       string           `json:"phone"`
	CountryCode string           `json:"country_code"`
	Status      string           `json:"status"`
	Bio         string           `json:"bio"`
	PublicUrl   string           `json:"public_url"`
	IsBlocked   bool             `json:"is_blocked"`
	IsApproved  bool             `json:"is_approved"`
	IsPublic    bool             `json:"is_public"`
	ReportCount int              `json:"report_count"`
	LastIp      string           `json:"last_ip"`
	LastLogin   time.Time        `json:"last_login"`
	Created     time.Time        `json:"created"`
	Updated     time.Time        `json:"updated"`
	Deleted     time.Time        `json:"deleted"`
	Role        RoleResponse     `json:"role,omitempty"`
	Reports     []ReportResponse `json:"reports,omitempty"`
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
		Bio:         user.Bio,
		PublicUrl:   user.PublicUrl,
		IsBlocked:   user.IsBlocked,
		IsApproved:  user.IsApproved,
		IsPublic:    user.IsPublic,
		ReportCount: user.ReportCount,
		LastIp:      user.LastIp,
		LastLogin:   user.LastLogin,
		Created:     user.Created,
		Updated:     user.Updated,
		Deleted:     user.Deleted,
	}

	if user.Role.ID != 0 {
		response.Role = NewRoleResponseFromDomainRole(user.Role)
	}

	// Map Reports to ReportResponse
	for _, report := range user.Reports {
		response.Reports = append(response.Reports, ReportResponse{
			Id:               report.ID,
			UserId:           report.UserId,
			ReportType:       report.ReportType,
			ReportReason:     report.ReportReason,
			ReportedByUserId: report.ReportedByUserId,
			ReportedAt:       report.ReportedAt,
		})
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

type ReportResponse struct {
	Id               uint      `json:"id"`
	UserId           uint      `json:"user_id"`
	ReportType       string    `json:"report_type"`
	ReportReason     string    `json:"report_reason"`
	ReportedByUserId uint      `json:"reported_by_user_id"`
	ReportedAt       time.Time `json:"reported_at"`
}
