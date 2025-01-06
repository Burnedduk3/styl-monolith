package errorhandler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ErrorResponse represents the response structure for sending error details in JSON format.
// Code is a string that identifies the specific error type.
// Message provides a human-readable description of the error.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HandleError processes and responds to HTTP errors based on domain-specific error codes or generic unexpected errors.
// It logs unexpected errors and maps domain error codes to appropriate HTTP status codes for client responses.
func HandleError(c echo.Context, err error, logger *logrus.Logger) error {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return c.JSON(getHTTPStatusCode(domainErr.Code), ErrorResponse{
			Code:    domainErr.Code,
			Message: domainErr.Message,
		})
	}
	logger.Error(err.Error())
	// Handle unexpected errorhandler
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "An unexpected error occurred",
	})
}

// getHTTPStatusCode maps specific error codes to their corresponding HTTP status codes.
// It returns 400 for recognized bad request error codes or 500 for unrecognized error codes.
func getHTTPStatusCode(errorCode string) int {
	switch errorCode {
	case ErrUserRequestPayloadBadRequestValidationError,
		ErrUserRequestPayloadBadRequestEmptyFields,
		ErrUserRequestPayloadBadRequestJsonBadlyFormated,
		ErrRoleRequestPayloadBadRequestJsonBadlyFormated,
		ErrRoleRequestPayloadBadRequestEmptyFields,
		ErrRoleRequestPayloadBadRequestValidationError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
