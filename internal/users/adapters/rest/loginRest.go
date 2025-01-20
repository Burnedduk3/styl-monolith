package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"styl-monolith/internal/users/adapters/rest/models"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/errorhandler"
	"styl-monolith/pkg/logger"
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
	crudService  services.CrudService
	logger       *logrus.Logger
}

func NewLoginRestUser(loginService services.LoginService, crudService services.CrudService, logger *logrus.Logger) LoginRest {
	return &LoginRestStruct{
		loginService: loginService,
		crudService:  crudService,
		logger:       logger,
	}
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

	userAuth, err := l.loginService.Login(body.Email, body.Password)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	l.logger.Debug("Fetching User from database with email: ", body.Email)
	dUser, err := l.crudService.GetUserByEmail(body.Email)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	loginRes := models.LoginRestResponse{
		User:      models.NewUserResponseFromDomainUser(dUser),
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
	userToBeCreated := body.ToUserDomain()
	l.logger.Info(fmt.Sprintf(logger.CreateUserOnDB, userToBeCreated.Email))
	err := l.loginService.SignUp(userToBeCreated)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	dUser, err := l.crudService.CreateUser(userToBeCreated)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	return ech.JSON(http.StatusOK, models.NewUserResponseFromDomainUser(dUser))
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
	userAuthResponse, err := l.loginService.Refresh(body.TokenId, body.Email)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}

	return ech.JSON(http.StatusOK, struct {
		TokenHash string `json:"tokenHash"`
	}{TokenHash: userAuthResponse.IdTokenHash})
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

	tokenHashI := ech.Get("tokenHash")
	tokenHash, ok := tokenHashI.(string)
	if !ok {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			"Invalid token hash format",
			nil)
		l.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, l.logger)
	}
	l.logger.Debug("Fetching User from database with email: ", body.Email)
	dUser, err := l.crudService.GetUserByEmail(body.Email)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	err = l.loginService.DeleteUserInCognito(dUser.Email, tokenHash)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	l.logger.Debug("Deleting User from database with id: ", dUser.ID)
	err = l.crudService.DeleteUserById(dUser.ID)
	if err != nil {
		l.logger.Error(err)
		return errorhandler.HandleError(ech, err, l.logger)
	}
	return ech.JSON(http.StatusOK, nil)
}
