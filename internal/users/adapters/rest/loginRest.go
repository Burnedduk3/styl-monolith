package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"styl-monolith/internal/users/adapters/rest/models"
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
	DeleteUserAccount(ech echo.Context) error
}

type LoginRestStruct struct {
	loginService services.LoginService
	logger       *logrus.Logger
}

func NewLoginRestUser(loginService services.LoginService, logger *logrus.Logger) LoginRest {
	return &LoginRestStruct{loginService: loginService, logger: logger}
}

func (l *LoginRestStruct) PostLogin(ech echo.Context) error {
	l.logger.Debug("Post login method called")
	var body models.LoginRestRequest
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		l.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, l.logger)
	}

	userAuth, user, err := l.loginService.Login(body.Email, body.Password)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	loginRes := models.LoginRestResponse{
		User:      models.NewUserResponseFromDomainUser(user),
		TokenHash: userAuth.IdTokenHash,
	}
	return ech.JSON(http.StatusOK, loginRes)
}

func (l *LoginRestStruct) PostSignUp(ech echo.Context) error {
	l.logger.Debug("Post signup method called")
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

func (l *LoginRestStruct) PostLogout(ech echo.Context) error {
	l.logger.Debug("Post refresh token method called")
	var body models.TokenRequest
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		l.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, l.logger)
	}
	err := l.loginService.Logout(body.TokenId, body.Email)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	return ech.JSON(http.StatusOK, nil)
}

func (l *LoginRestStruct) PostRefreshToken(ech echo.Context) error {
	l.logger.Debug("Post refresh token method called")
	var body models.TokenRequest
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		l.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, l.logger)
	}
	dUser, err := l.loginService.Refresh(body.TokenId, body.Email)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	return ech.JSON(http.StatusOK, dUser)
}

func (l *LoginRestStruct) PatchChangePassword(ech echo.Context) error {
	// TODO: Implement UpdatePassword
	return nil
}

func (l *LoginRestStruct) DeleteUserAccount(ech echo.Context) error {
	l.logger.Debug("Delete User account method called")
	var body models.TokenRequest
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		l.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, l.logger)
	}
	err := l.loginService.DeleteUser(body.Email)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	return ech.JSON(http.StatusOK, nil)
}
