package errorhandler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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

func getHTTPStatusCode(errorCode string) int {
	switch errorCode {
	case ErrUserRequestPayloadBadRequestValidationError,
		ErrUserRequestPayloadBadRequestEmptyFields,
		ErrUserRequestPayloadBadRequestJsonBadlyFormated:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
