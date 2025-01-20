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
	ReportCount int
	Username    string
	Name        string
	LastName    string
	Birthday    string
	Email       string
	Country     string
	Phone       string
	CountryCode string
	Password    string
	Bio         string
	PublicUrl   string
	LastIp      string
	IsBlocked   bool
	IsPublic    bool
	IsApproved  bool
	Role        Role
	Status      Status
	Reports     []Report
	LastLogin   time.Time
	Created     time.Time
	Updated     time.Time
	Deleted     time.Time
}

func (u User) UserToUserPayload() UserPayload {
	return UserPayload{
		Name:        u.Name,
		LastName:    u.LastName,
		Username:    u.Username,
		Birthday:    u.Birthday,
		Email:       u.Email,
		Country:     u.Country,
		Password:    u.Password,
		IsBlocked:   u.IsBlocked,
		IsPublic:    u.IsPublic,
		PublicUrl:   u.PublicUrl,
		IsApproved:  u.IsApproved,
		ReportCount: u.ReportCount,
		Bio:         u.Bio,
		CountryCode: u.CountryCode,
		Phone:       u.Phone,
	}
}

type UserPayload struct {
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Birthday    string `json:"birthday"`
	Email       string `json:"email"`
	Country     string `json:"country"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	CountryCode string `json:"country_code"`
	IsBlocked   bool   `json:"is_blocked"`
	IsPublic    bool   `json:"is_public"`
	PublicUrl   string `json:"public_url"`
	IsApproved  bool   `json:"is_approved"`
	ReportCount int    `json:"report_count"`
	Bio         string `json:"bio"`
}

func (c *UserPayload) ToUserDomain() User {
	return User{
		Name:        c.Name,
		LastName:    c.LastName,
		Username:    c.Username,
		Birthday:    c.Birthday,
		Email:       c.Email,
		Country:     c.Country,
		CountryCode: c.CountryCode,
		Password:    c.Password,
		Phone:       c.Phone,
		IsBlocked:   c.IsBlocked,
		IsPublic:    c.IsPublic,
		PublicUrl:   c.PublicUrl,
		IsApproved:  c.IsApproved,
		ReportCount: c.ReportCount,
		Bio:         c.Bio,
	}
}

func (c *UserPayload) Validate() error {
	if validator.IsStringEmpty(c.Name) ||
		validator.IsStringEmpty(c.LastName) ||
		validator.IsStringEmpty(c.Birthday) ||
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
	isValidBio := validator.ValidateString(c.Bio)
	if !isValidContryCode ||
		!isValidPhone ||
		!isValidBio ||
		!isValidEmail {
		return errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestValidationError),
			nil)
	}
	return nil
}

type Report struct {
	ID               uint
	UserId           uint
	ReportType       string
	ReportReason     string
	ReportedByUserId uint
	ReportedAt       time.Time
}
