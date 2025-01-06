package errorhandler

const (
	// User domain errorhandler
	ErrUserRequestPayloadBadRequestValidationError   = "USER_001"
	ErrUserRequestPayloadBadRequestEmptyFields       = "USER_002"
	ErrUserRequestPayloadBadRequestJsonBadlyFormated = "USER_003"
	ErrUserDatabaseUnableToCompleteOperation         = "USER_004"
	ErrUserNotFound                                  = "USER_005"
	ErrUserIncompleteParamsForCompleteUpdate         = "USER_006"

	// Role error handler
	ErrRoleRequestPayloadBadRequestJsonBadlyFormated = "ROLE_001"
	ErrRoleRequestPayloadBadRequestEmptyFields       = "ROLE_002"
	ErrRoleDatabaseUnableToCompleteOperation         = "ROLE_003"
	ErrRoleNotFound                                  = "ROLE_004"
	ErrRoleIncompleteParamsForCompleteUpdate         = "ROLE_005"
	ErrRoleRequestPayloadBadRequestValidationError   = "ROLE_006"
)

// errorMessages defines a mapping of error codes to their corresponding descriptive error messages.
var errorMessages = map[string]string{
	ErrUserRequestPayloadBadRequestValidationError:   "request payload is wrong, had validation errors on inputs",
	ErrUserRequestPayloadBadRequestEmptyFields:       "request payload is wrong, Required fields are empty",
	ErrUserRequestPayloadBadRequestJsonBadlyFormated: "request payload for user is wrong, JSON is badly formated",
	ErrRoleRequestPayloadBadRequestJsonBadlyFormated: "request payload for role is wrong, JSON is badly formated",
	ErrRoleRequestPayloadBadRequestEmptyFields:       "request payload for role is wrong, required fields are empty",
	ErrRoleDatabaseUnableToCompleteOperation:         "unable to complete operation on specified role",
	ErrUserDatabaseUnableToCompleteOperation:         "unable to complete operation on specified user",
	ErrRoleRequestPayloadBadRequestValidationError:   "request payload for role is wrong, validation errors on inputs",
	ErrRoleNotFound:                          "There is no role with the specified param: %v",
	ErrUserNotFound:                          "There is no user with the specified param: %v",
	ErrRoleIncompleteParamsForCompleteUpdate: "Incomplete role params for complete update",
	ErrUserIncompleteParamsForCompleteUpdate: "Incomplete user params for complete update",
}

// GetErrorMessage returns a human-readable error message corresponding to the provided error code.
// If the error code does not match any predefined messages, it returns "unknown error".
func GetErrorMessage(code string) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "unknown error"
}
