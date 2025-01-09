package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/errorhandler"
)

type LoginRest interface {
	PostLogin(ech echo.Context) error
	PostSignUp(ech echo.Context) error
	PostLogout(ech echo.Context) error
	PostRefreshToken(ech echo.Context) error
	PatchChangePassword(ech echo.Context) error
}

type LoginRestStruct struct {
	loginService services.LoginService
	logger       *logrus.Logger
}

func NewLoginRestUser(loginService services.LoginService, logger *logrus.Logger) LoginRest {
	return &LoginRestStruct{loginService: loginService, logger: logger}
}

func (l *LoginRestStruct) PostLogin(c echo.Context) error {
	// TODO: Implement PostLogin
	return nil
}

func (l *LoginRestStruct) PostSignUp(ech echo.Context) error {
	l.logger.Debug("PatchUser method called")
	var body domain.UserPayload
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		l.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, l.logger)
	}
	if err := body.Validate(); err != nil {
		l.logger.Error(err.Error())
		return errorhandler.HandleError(ech, err, l.logger)
	}
	dUser, err := l.loginService.SignUp(body.ToUserDomain())
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	return ech.JSON(http.StatusOK, dUser)
}

func (l *LoginRestStruct) PostLogout(c echo.Context) error {
	// TODO: Implement PostLogout
	return nil
}

func (l *LoginRestStruct) PostRefreshToken(c echo.Context) error {
	// TODO: Implement PostRefresh
	return nil
}

func (l *LoginRestStruct) PatchChangePassword(c echo.Context) error {
	// TODO: Implement UpdatePassword
	return nil
}
