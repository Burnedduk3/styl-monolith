package models

import (
	"gorm.io/gorm"
	"styl-monolith/internal/users/core/domain"
	"time"
)

type User struct {
	gorm.Model
	Name      string
	Email     string
	Country   string
	RoleID    uint
	Status    string `gorm:"default:'active'"`
	LastIp    *string
	LastLogin time.Time `gorm:"autoUpdateTime:milli"`
}

func NewPostgresUserFromDomainUser(user domain.User) User {
	return User{
		Name:      user.Name,
		Email:     user.Email,
		Country:   user.Country,
		Status:    string(user.Status),
		LastIp:    &user.LastIp,
		LastLogin: user.LastLogin,
	}
}

func (u *User) toUserDomain() domain.User {
	mappedUser := domain.User{
		ID:        u.Model.ID,
		Name:      u.Name,
		Email:     u.Email,
		Country:   u.Country,
		LastLogin: u.LastLogin,
		Created:   u.CreatedAt,
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
