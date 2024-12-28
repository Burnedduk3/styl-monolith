package errorhandler

const (
	// User domain errorhandler
	ErrUserRequestPayloadBadRequestValidationError   = "USER_001"
	ErrUserRequestPayloadBadRequestEmptyFields       = "USER_002"
	ErrUserRequestPayloadBadRequestJsonBadlyFormated = "USER_003"
	ErrUserDatabaseUnableToCompleteOperation         = "USER_004"

	// Role error handler
	ErrRoleRequestPayloadBadRequestJsonBadlyFormated = "ROLE_001"
	ErrRoleRequestPayloadBadRequestEmptyFields       = "ROLE_002"
	ErrRoleDatabaseUnableToCompleteOperation         = "ROLE_003"
)

var errorMessages = map[string]string{
	ErrUserRequestPayloadBadRequestValidationError:   "request payload is wrong, had validation errors on inputs",
	ErrUserRequestPayloadBadRequestEmptyFields:       "request payload is wrong, Required fields are empty",
	ErrUserRequestPayloadBadRequestJsonBadlyFormated: "request payload for user is wrong, JSON is badly formated",
	ErrRoleRequestPayloadBadRequestJsonBadlyFormated: "request payload for role is wrong, JSON is badly formated",
	ErrRoleRequestPayloadBadRequestEmptyFields:       "request payload for role is wrong, required fields are empty",
	ErrRoleDatabaseUnableToCompleteOperation:         "unable to complete operation on specified role",
	ErrUserDatabaseUnableToCompleteOperation:         "unable to complete operation on specified user",
}

func GetErrorMessage(code string) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "unknown error"
}
