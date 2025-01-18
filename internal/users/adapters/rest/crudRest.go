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

// CrudRest defines methods for CRUD operations on users and roles in a RESTful API.
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

// crudRestStruct represents a REST handler structure for CRUD operations on users and roles.
// It utilizes an injected CrudService for handling business logic.
// Logger is used for structured logging within the CRUD operations.
type crudRestStruct struct {
	crudService services.CrudService
	logger      *logrus.Logger
}

// NewCrudRestUser initializes and returns a new instance of CrudRest using the provided CrudService and logger.
func NewCrudRestUser(crudService services.CrudService, logger *logrus.Logger) CrudRest {
	return &crudRestStruct{crudService: crudService, logger: logger}
}

// CreateUser handles the creation of a new user from the request payload and returns the created user as a response.
// It validates the request payload, processes it using the service layer, and handles errors appropriately.
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

// CreateRole handles the creation of a new role based on the provided payload and returns the created role as a response.
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

// DeleteUserById deletes a user identified by their string-based ID from the system.
// It parses the string ID to uint, validates it, and calls the service layer for the deletion operation.
// Returns an error if the ID cannot be parsed, the deletion fails, or if any validation error occurs.
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

// DeleteRoleById deletes a role by its ID provided as a string, returning an error if the operation fails.
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

// UpdateUserById updates a user by their ID using data from the request payload and returns the updated user or an error.
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

// UpdateRoleById updates a role in the system based on the provided ID and request payload.
// It validates the input, binds the request payload to a RolePayload struct, and processes the update.
// Returns an error if validation, binding, or the update operation fails.
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

// PatchUserById partially updates a user's information by their ID, based on the provided payload in the request context.
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

// PatchRoleById partially updates a role identified by stringId using the received JSON payload and returns the updated role.
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

// GetUserById retrieves a user by their ID, parses the ID from a string, and returns the user data in JSON format.
// If the ID is invalid or an error occurs, it handles the error and sends an appropriate response.
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

// GetUserByQuery retrieves a user based on query parameters (phone, username, or email) and returns the user data as JSON.
// It handles potential errors during the process and logs them.
func (c *crudRestStruct) GetUserByQuery(ech echo.Context) error {
	c.logger.Debug("GetUserByQuery method called")
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

// GetRoleById retrieves a role by its unique ID, parses it from string to uint, and returns it as a JSON response.
// Handles any errors during ID conversion or service layer interaction, responding with appropriate error messages.
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

// ListUsers handles the HTTP request to fetch a paginated list of users by processing query parameters like page and size.
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

// ListRoles retrieves a paginated list of roles from the service layer and returns them in a JSON response format.
// It parses query parameters for pagination (`page` and `size`) and ensures default values if parameters are invalid.
// The response includes roles mapped to a RoleResponse model and pagination metadata.
// Returns an error response if role retrieval fails.
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

// UpdateUserRoleById updates the role of a user by their ID using the provided payload and returns the updated user information.
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
