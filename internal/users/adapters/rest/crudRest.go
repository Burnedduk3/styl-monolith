package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	crudService services.CrudService
	logger      *logrus.Logger
}

func NewCrudRestUser(crudService services.CrudService, logger *logrus.Logger) CrudRest {
	return &crudRestStruct{crudService: crudService, logger: logger}
}

func (c *crudRestStruct) CreateUser(ech echo.Context) error {
	c.logger.Debug("CreateUser method called")
	var body models.UserPayload
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	if err := body.Validate(); err != nil {
		c.logger.Error(err.Error())
		return errorhandler.HandleError(ech, err, c.logger)
	}
	dUser, err := c.crudService.CreateUser(body.Name, body.Username, body.Email, body.Country, body.Phone, body.CountryCode)
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewUserResponseFromDomainUser(dUser))
}

func (c *crudRestStruct) CreateRole(ech echo.Context) error {
	c.logger.Debug("CreateRole method called")
	var body models.RolePayload
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	if err := body.Validate(); err != nil {
		c.logger.Error(err.Error())
		return errorhandler.HandleError(ech, err, c.logger)
	}
	dRole, err := c.crudService.CreateRole(body.Name, body.Description)
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewRoleResponseFromDomainRole(dRole))
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
	// Parse query parameters
	pageParam := ech.QueryParam("page")
	sizeParam := ech.QueryParam("size")

	// Default values
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil || size < 1 {
		size = 10
	}
	var users []models.UserResponse
	dUsers, totalPages, err := c.crudService.ListUsers(page, size)
	for _, dUser := range dUsers {
		users = append(users, models.NewUserResponseFromDomainUser(dUser))
	}
	pagResponse := models.PaginationResponse{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  totalPages,
		Data:        users,
	}
	return ech.JSON(http.StatusOK, pagResponse)
}

func (c *crudRestStruct) ListRoles(ech echo.Context) error {
	c.logger.Debug("ListRoles method called")
	// Parse query parameters
	pageParam := ech.QueryParam("page")
	sizeParam := ech.QueryParam("size")

	// Default values
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil || size < 1 {
		size = 10
	}
	var roles []models.RoleResponse
	dRoles, totalPages, err := c.crudService.ListRoles(page, size)
	for _, dRole := range dRoles {
		roles = append(roles, models.NewRoleResponseFromDomainRole(dRole))
	}
	pagResponse := models.PaginationResponse{
		CurrentPage: page,
		PageSize:    size,
		TotalPages:  totalPages,
		Data:        roles,
	}
	return ech.JSON(http.StatusOK, pagResponse)
}
