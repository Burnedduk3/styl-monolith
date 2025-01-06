package services

import (
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/ports"
)

type LoginServiceStruct struct {
	logger  *logrus.Logger
	awsAuth ports.LoginPort
}

type LoginService interface {
	Login() error
	Logout() error
	Refresh() error
	ChangePassword() error
}

func NewLoginService(log *logrus.Logger, loginPort ports.LoginPort) LoginService {
	return &LoginServiceStruct{
		logger:  log,
		awsAuth: loginPort,
	}
}

func (l *LoginServiceStruct) Login() error {
	// TODO: Implement Login logic
	return nil
}

func (l *LoginServiceStruct) Logout() error {
	// TODO: Implement Logout logic
	return nil
}

func (l *LoginServiceStruct) Refresh() error {
	// TODO: Implement Refresh logic
	return nil
}

func (l *LoginServiceStruct) ChangePassword() error {
	// TODO: Implement ChangePassword logic
	return nil
}
