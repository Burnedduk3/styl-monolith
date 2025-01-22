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

// CrudRest defines an interface for CRUD operations on users, roles, and reports in a REST API context.
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
	CreateReport(ech echo.Context) error
	DeleteReportById(ech echo.Context, id string) error
	GetReportById(ech echo.Context, id string) error
	GetUserReportsByUserId(ech echo.Context, userId string) error
}

// crudRestStruct is a struct that handles operations for a CRUD service and provides logging functionality.
type crudRestStruct struct {
	crudService services.CrudService
	logger      *logrus.Logger
}

// NewCrudRestUser creates a new CrudRest implementation with the provided CrudService and logger.
func NewCrudRestUser(crudService services.CrudService, logger *logrus.Logger) CrudRest {
	return &crudRestStruct{crudService: crudService, logger: logger}
}

// CreateUser processes a HTTP request to create a new user based on the provided payload and returns a JSON response.
// It validates the payload, converts it to a domain object, and interacts with the service layer for user creation.
// In case of errors, it handles them appropriately and sends the corresponding HTTP response.
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

// CreateRole handles the creation of a new role based on the provided payload, validates it, and returns the created role.
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

// DeleteUserById deletes a user identified by the given string ID, handling validation and errors during the process.
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

// DeleteRoleById deletes a role identified by its ID, parses the string ID to uint64, and handles errors during the process.
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

// UpdateUserById handles the HTTP request to update a user by their ID based on the provided payload.
// It validates the input, processes the update via the service layer, and returns the updated user in the response.
// If any errors occur, appropriate domain or validation errors are logged and handled.
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

// UpdateRoleById handles the update of a role identified by the given ID using the supplied payload and request context.
// It validates the payload, converts it to the domain model, and invokes the service layer for persistence.
// Returns HTTP 400 for validation or binding errors and HTTP 500 for internal server errors.
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

// PatchUserById partially updates user data based on the given user ID and request payload, responding with the updated user.
// It handles JSON binding errors, input validation, and domain-specific errors, returning appropriate HTTP responses.
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

// PatchRoleById updates a role by its ID using the provided partial update payload.
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

// GetUserById retrieves a user by their ID, validates the input, and returns the user details in JSON format.
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

// GetUserByQuery retrieves a user based on query parameters (phone, username, email) from the HTTP request context.
// It uses the appropriate service methods depending on the query parameters provided and returns the user as JSON.
// If an error occurs during retrieval, it is logged and appropriately handled with an HTTP error response.
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

// GetRoleById retrieves a role by its ID, validates the input, and returns the role as a JSON response or an error response.
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

// ListUsers retrieves a paginated list of users based on query parameters for page and size.
// It defaults to page 1 and size 10 if parameters are invalid or not provided.
// Returns a JSON response containing user data and pagination metadata.
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

// ListRoles retrieves a paginated list of roles based on query parameters (page and size) and returns them in JSON format.
// It validates input parameters, handles pagination defaults, maps domain roles to RoleResponse, and logs errors or events.
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

// UpdateUserRoleById updates the role of a user identified by their ID, using the provided role ID in the request payload.
// It validates the request payload and the user ID, then performs the update operation via the CRUD service.
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

// CreateReport handles the creation of a new report, including validation, logging, and communication with the CRUD service.
// It binds the request payload, validates it, processes domain conversion, and responds with the created report or errors.
// Returns a JSON response with an HTTP status, either `http.StatusCreated` or an appropriate error status code.
func (c *crudRestStruct) CreateReport(ech echo.Context) error {
	c.logger.Debug("CreateReport method called")

	var body models.ReportPayload
	if err := ech.Bind(&body); err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrReportRequestPayloadBadRequestJsonBadlyFormated,
			errorhandler.GetErrorMessage(errorhandler.ErrReportRequestPayloadBadRequestJsonBadlyFormated),
			err,
		)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}

	if err := body.Validate(); err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}

	report := body.ToReportDomain()
	dReport, err := c.crudService.CreateUserReport(domain.User{ID: report.UserId}, report)
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}

	return ech.JSON(http.StatusCreated, models.NewReportResponseFromDomainReport(dReport))
}

// DeleteReportById removes a report by its ID, validates the input, and handles errors using the appropriate error handler.
func (c *crudRestStruct) DeleteReportById(ech echo.Context, stringId string) error {
	c.logger.Debug("DeleteReportById method called")

	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrReportRequestPayloadValidationFailed,
			errorhandler.GetErrorMessage(errorhandler.ErrReportRequestPayloadValidationFailed),
			err,
		)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}

	report := domain.Report{ID: uint(id)}
	_, err = c.crudService.DeleteUserReport(domain.User{ID: report.UserId}, report)
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}

	return ech.JSON(http.StatusOK, nil)
}

// GetReportById retrieves a report by its ID, validates the ID input, and returns the report in the response or an error.
func (c *crudRestStruct) GetReportById(ech echo.Context, stringId string) error {
	c.logger.Debug("GetReportById method called")

	id, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrReportRequestPayloadValidationFailed,
			errorhandler.GetErrorMessage(errorhandler.ErrReportRequestPayloadValidationFailed),
			err,
		)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}

	report, err := c.crudService.GetReportById(uint(id))
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}

	return ech.JSON(http.StatusOK, models.NewReportResponseFromDomainReport(report))
}

// GetUserReportsByUserId retrieves user reports based on the specified user ID and returns the result as JSON response.
// It validates the user ID, handles potential errors, and logs relevant actions during the request processing.
func (c *crudRestStruct) GetUserReportsByUserId(ech echo.Context, userId string) error {
	c.logger.Debug("GetUserReports method called")

	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		domErr := errorhandler.NewDomainError(
			errorhandler.ErrUserRequestPayloadBadRequestValidationError,
			errorhandler.GetErrorMessage(errorhandler.ErrUserRequestPayloadBadRequestValidationError),
			err,
		)
		c.logger.Error(domErr.Error())
		return errorhandler.HandleError(ech, domErr, c.logger)
	}

	reports, err := c.crudService.GetUserReports(domain.User{ID: uint(id)})
	if err != nil {
		c.logger.Error(err)
		return errorhandler.HandleError(ech, err, c.logger)
	}

	var response []models.ReportResponse
	for _, report := range reports {
		response = append(response, models.NewReportResponseFromDomainReport(report))
	}

	return ech.JSON(http.StatusOK, response)
}
