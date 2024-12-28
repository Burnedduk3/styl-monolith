package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"styl-monolith/internal/users/core/services"
)

type LoginRest interface {
	PostLogin(ech echo.Context) error
	PostLogout(ech echo.Context) error
	PostRefreshToken(ech echo.Context) error
	ChangePassword(ech echo.Context) error
}

type LoginRestStruct struct {
	loginService *services.LoginService
	logger       *logrus.Logger
}

func NewLoginRestUser(loginService *services.LoginService, logger *logrus.Logger) LoginRest {
	return &LoginRestStruct{loginService: loginService, logger: logger}
}

func (l *LoginRestStruct) PostLogin(c echo.Context) error {
	// TODO: Implement PostLogin
	return nil
}

func (l *LoginRestStruct) PostLogout(c echo.Context) error {
	// TODO: Implement PostLogout
	return nil
}

func (l *LoginRestStruct) PostRefreshToken(c echo.Context) error {
	// TODO: Implement PostRefresh
	return nil
}

func (l *LoginRestStruct) ChangePassword(c echo.Context) error {
	// TODO: Implement UpdatePassword
	return nil
}
