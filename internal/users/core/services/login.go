package services

import (
	"github.com/sirupsen/logrus"
)

type LoginServiceStruct struct {
	logger *logrus.Logger
}

type LoginService interface {
	Login() error
	Logout() error
	Refresh() error
	ChangePassword() error
}

func NewLoginService(log *logrus.Logger) LoginService {
	return &LoginServiceStruct{
		logger: log,
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
