package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"styl-monolith/internal/users/adapters/rest/models"
	"styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/errorhandler"
)

type CrudRest interface {
	CreateUser(ech echo.Context) error
	CreateRole(ech echo.Context) error
	DeleteUser(ech echo.Context, id uint) error
	DeleteRole(ech echo.Context, id uint) error
	UpdateUser(ech echo.Context, id uint) error
	UpdateRole(ech echo.Context, id uint) error
	GetUser(ech echo.Context, id uint) error
	GetRole(ech echo.Context, id uint) error
	ListUsers(ech echo.Context) error
	ListRoles(ech echo.Context) error
}

type crudRestStruct struct {
	crudService *services.CrudService
	logger      *logrus.Logger
}

func NewCrudRestUser(crudService *services.CrudService, logger *logrus.Logger) CrudRest {
	return &crudRestStruct{crudService: crudService, logger: logger}
}

func (c *crudRestStruct) CreateUser(ech echo.Context) error {
	c.logger.Debug("CreateUser method called")
	var body models.CreateUserPayload
	if err := ech.Bind(&body); err != nil || body.Validate() != nil {
		domErr := errorhandler.NewDomainError(errorhandler.ErrUserRequestPayloadBadRequest, errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequest), err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}

	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) CreateRole(ech echo.Context) error {
	c.logger.Debug("CreateRole method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) DeleteUser(ech echo.Context, id uint) error {
	c.logger.Debug("DeleteUser method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) DeleteRole(ech echo.Context, id uint) error {
	c.logger.Debug("DeleteRole method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) UpdateUser(ech echo.Context, id uint) error {
	c.logger.Debug("UpdateUser method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) UpdateRole(ech echo.Context, id uint) error {
	c.logger.Debug("UpdateRole method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) GetUser(ech echo.Context, id uint) error {
	c.logger.Debug("GetUser method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) GetRole(ech echo.Context, id uint) error {
	c.logger.Debug("GetRole method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) ListUsers(ech echo.Context) error {
	c.logger.Debug("ListUsers method called")
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) ListRoles(ech echo.Context) error {
	c.logger.Debug("ListRoles method called")
	return ech.JSON(http.StatusOK, nil)
}
