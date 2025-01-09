package models

import "styl-monolith/internal/users/core/domain"

type CognitoUser struct {
	ID       uint
	Name     string
	LastName string
	Email    string
	Password string
	Username string
	Phone    string
}

func NewCognitoUserFromDomain(user domain.User) CognitoUser {
	return CognitoUser{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		Phone:    user.Phone,
	}
}

func (c *CognitoUser) NewDomainUserFromCognito() domain.User {
	return domain.User{
		ID:       c.ID,
		Name:     c.Name,
		LastName: c.LastName,
		Email:    c.Email,
		Username: c.Username,
		Password: c.Password,
		Phone:    c.Phone,
	}
}
