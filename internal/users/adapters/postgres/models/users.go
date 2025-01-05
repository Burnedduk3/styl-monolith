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
	Email       string `gorm:"unique"`
	Phone       string `gorm:"unique"`
	CountryCode string
	Country     string
	RoleID      uint
	Status      string `gorm:"default:'active'"`
	LastIp      *string
	LastLogin   time.Time `gorm:"autoUpdateTime:milli"`
}

func NewPostgresUserFromDomainUser(user domain.User) User {
	u := User{
		Model: gorm.Model{
			ID: user.ID,
		},
		Username:    user.Username,
		Name:        user.Name,
		Phone:       user.Phone,
		CountryCode: user.CountryCode,
		Email:       user.Email,
		Country:     user.Country,
		Status:      string(user.Status),
		LastIp:      &user.LastIp,
		LastLogin:   user.LastLogin,
	}
	if user.Role.ID != 0 {
		u.RoleID = user.Role.ID
	}
	return u
}

func (u *User) ToUserDomain() domain.User {
	mappedUser := domain.User{
		ID:       u.Model.ID,
		Username: u.Username,
		Name:     u.Name,
		Email:    u.Email,
		Phone:    u.Phone,
		Country:  u.Country,
		Role: domain.Role{
			ID: u.RoleID,
		},
		CountryCode: u.CountryCode,
		LastLogin:   u.LastLogin,
		Created:     u.CreatedAt,
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
