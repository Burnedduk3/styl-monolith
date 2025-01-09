package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"styl-monolith/internal/users/adapters/rest/models"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/services"
	"styl-monolith/pkg/errorhandler"
)

type CrudRest interface {
	CreateUser(ech echo.Context) error
	CreateRole(ech echo.Context) error
	DeleteUserById(ech echo.Context, id string) error
	DeleteRoleById(ech echo.Context, id string) error
	UpdateUserById(ech echo.Context, id string) error
	UpdateRoleById(ech echo.Context, id string) error
	PatchUserById(ech echo.Context, id string) error
	PatchRoleById(ech echo.Context, id string) error
	GetUserById(ech echo.Context, id string) error
	GetRoleById(ech echo.Context, id string) error
	GetUserByQuery(ech echo.Context) error
	ListUsers(ech echo.Context) error
	ListRoles(ech echo.Context) error
	UpdateUserRoleById(ech echo.Context, id string) error
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
	var body domain.UserPayload
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
	dUser, err := c.crudService.CreateUser(body.ToUserDomain())
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
	dRole, err := c.crudService.CreateRole(body.ToRoleDomain())
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewRoleResponseFromDomainRole(dRole))
}

func (c *crudRestStruct) DeleteUserById(ech echo.Context, stringId string) error {
	c.logger.Debug("DeleteUser method called")
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestValidationError),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	err = c.crudService.DeleteUserById(uint(id))
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) DeleteRoleById(ech echo.Context, stringId string) error {
	c.logger.Debug("DeleteRole method called")
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	err = c.crudService.DeleteRoleById(uint(id))
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, nil)
}

func (c *crudRestStruct) UpdateUserById(ech echo.Context, stringId string) error {
	c.logger.Debug("UpdateUser method called")
	var body domain.UserPayload
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
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	dUser, err := c.crudService.UpdateUserById(uint(id), body.ToUserDomain())
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, dUser)
}

func (c *crudRestStruct) UpdateRoleById(ech echo.Context, stringId string) error {
	c.logger.Debug("UpdateRole method called")
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
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	dRole, err := c.crudService.UpdateRoleById(uint(id), body.ToRoleDomain())
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewRoleResponseFromDomainRole(dRole))
}

func (c *crudRestStruct) PatchUserById(ech echo.Context, stringId string) error {
	c.logger.Debug("PatchUser method called")
	var body domain.UserPayload
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	dUser, err := c.crudService.PartialUpdateUserById(uint(id), body.ToUserDomain())
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, dUser)
}

func (c *crudRestStruct) PatchRoleById(ech echo.Context, stringId string) error {
	c.logger.Debug("UpdateRole method called")
	var body models.RolePayload
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	dRole, err := c.crudService.PartialUpdateRoleById(uint(id), body.ToRoleDomain())
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewRoleResponseFromDomainRole(dRole))
}

func (c *crudRestStruct) GetUserById(ech echo.Context, stringId string) error {
	c.logger.Debug("GetUserById method called")
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	dUser, err := c.crudService.GetUserById(uint(id))
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewUserResponseFromDomainUser(dUser))
}

func (c *crudRestStruct) GetUserByQuery(ech echo.Context) error {
	c.logger.Debug("GetUserByPhone method called")
	var dUser domain.User
	var err error
	phone := ech.QueryParam("phone")
	if phone != "" {
		dUser, err = c.crudService.GetUserByPhone(phone)
	}
	username := ech.QueryParam("username")
	if username != "" {
		dUser, err = c.crudService.GetUserByUsername(username)
	}
	email := ech.QueryParam("email")
	if email != "" {
		dUser, err = c.crudService.GetUserByEmail(email)
	}
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewUserResponseFromDomainUser(dUser))
}

func (c *crudRestStruct) GetRoleById(ech echo.Context, stringId string) error {
	c.logger.Debug("GetRole method called")
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	role, err := c.crudService.GetRoleById(uint(id))
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, models.NewRoleResponseFromDomainRole(role))
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
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
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
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
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

func (c *crudRestStruct) UpdateUserRoleById(ech echo.Context, stringId string) error {
	c.logger.Debug("Update User role method called")
	var body models.UpdateUserRolePayload
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrRoleRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrRoleRequestPayloadBadRequestJsonBadlyFormated),
			err)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}
	dUser, err := c.crudService.UpdateUserRole(uint(id), body.RoleId)
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}
	return ech.JSON(http.StatusOK, dUser)
}
