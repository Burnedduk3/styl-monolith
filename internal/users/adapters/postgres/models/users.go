package models

import (
	"gorm.io/gorm"
	"styl-monolith/internal/users/core/domain"
	"time"
)

type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Name        string
	LastName    string
	Birthday    string
	Email       string `gorm:"unique"`
	Phone       string `gorm:"unique"`
	CountryCode string
	Country     string
	RoleID      uint
	Reports     []Report `gorm:"foreignKey:UserId"`
	IsBlocked   bool     `gorm:"default:false"`
	IsApproved  bool     `gorm:"default:false"`
	IsPublic    bool     `gorm:"default:true"`
	Bio         string   `gorm:"type:text"`
	PublicUrl   string
	ReportCount int    `gorm:"default:0"`
	Status      string `gorm:"default:'active'"`
	LastIp      *string
	LastLogin   time.Time `gorm:"autoUpdateTime:milli"`
}

type Report struct {
	gorm.Model
	UserId           uint `gorm:"index"`
	ReportType       string
	ReportReason     string
	ReportedByUserId uint
	ReportedAt       time.Time
}

func NewPostgresUserFromDomainUser(user domain.User) User {
	u := User{
		Model: gorm.Model{
			ID: user.ID,
		},
		Username:    user.Username,
		Name:        user.Name,
		LastName:    user.LastName,
		Birthday:    user.Birthday,
		Phone:       user.Phone,
		CountryCode: user.CountryCode,
		Email:       user.Email,
		Country:     user.Country,
		IsBlocked:   user.IsBlocked,
		IsApproved:  user.IsApproved,
		IsPublic:    user.IsPublic,
		Bio:         user.Bio,
		PublicUrl:   user.PublicUrl,
		ReportCount: user.ReportCount,
		Reports:     make([]Report, 0),
		Status:      string(user.Status),
		LastIp:      &user.LastIp,
		LastLogin:   user.LastLogin,
	}
	if user.Role.ID != 0 {
		u.RoleID = user.Role.ID
	}
	// Map Reports from the domain user
	for _, report := range user.Reports {
		u.Reports = append(u.Reports, Report{
			Model: gorm.Model{
				ID: report.ID,
			},
			UserId:           report.UserId,
			ReportType:       report.ReportType,
			ReportReason:     report.ReportReason,
			ReportedByUserId: report.ReportedByUserId,
			ReportedAt:       report.ReportedAt,
		})
	}
	return u
}

func (u *User) ToUserDomain() domain.User {
	mappedUser := domain.User{
		ID:       u.Model.ID,
		Username: u.Username,
		Name:     u.Name,
		LastName: u.LastName,
		Birthday: u.Birthday,
		Email:    u.Email,
		Phone:    u.Phone,
		Country:  u.Country,
		Role: domain.Role{
			ID: u.RoleID,
		},
		CountryCode: u.CountryCode,
		IsBlocked:   u.IsBlocked,
		IsApproved:  u.IsApproved,
		IsPublic:    u.IsPublic,
		Bio:         u.Bio,
		PublicUrl:   u.PublicUrl,
		ReportCount: u.ReportCount,
		LastLogin:   u.LastLogin,
		Created:     u.CreatedAt,
	}

	// Map Reports to the domain user
	for _, report := range u.Reports {
		mappedUser.Reports = append(mappedUser.Reports, domain.Report{
			ID:               report.ID,
			UserId:           report.UserId,
			ReportType:       report.ReportType,
			ReportReason:     report.ReportReason,
			ReportedByUserId: report.ReportedByUserId,
			ReportedAt:       report.ReportedAt,
		})
	}

	if u.Status == "active" {
		mappedUser.Status = domain.StatusActive
	} else if u.Status == "inactive" {
		mappedUser.Status = domain.StatusInactive
	} else if u.Status == "deleted" {
		mappedUser.Status = domain.StatusDeleted
	}
	if u.LastIp != nil {
		mappedUser.LastIp = *u.LastIp
	}
	return mappedUser
}
