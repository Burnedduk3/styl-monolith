package errorhandler

const (
	// User domain errors
	ErrUserRequestPayloadBadRequestValidationError   = "USER_001"
	ErrUserRequestPayloadBadRequestEmptyFields       = "USER_002"
	ErrUserRequestPayloadBadRequestJsonBadlyFormated = "USER_003"
	ErrUserDatabaseUnableToCompleteOperation         = "USER_004"
	ErrUserNotFound                                  = "USER_005"
	ErrUserIncompleteParamsForCompleteUpdate         = "USER_006"
	ErrCreatingUserInCognitoUserPool                 = "USER_007"
	ErrSettingPermanentPassword                      = "USER_008"
	ErrDeletingUserFromCognitoUserPool               = "USER_009"

	// Role domain errors
	ErrRoleRequestPayloadBadRequestJsonBadlyFormated = "ROLE_001"
	ErrRoleRequestPayloadBadRequestEmptyFields       = "ROLE_002"
	ErrRoleDatabaseUnableToCompleteOperation         = "ROLE_003"
	ErrRoleNotFound                                  = "ROLE_004"
	ErrRoleIncompleteParamsForCompleteUpdate         = "ROLE_005"
	ErrRoleRequestPayloadBadRequestValidationError   = "ROLE_006"

	// User report errors
	ErrReportDatabaseUnableToCompleteOperation         = "USER_REPORT_001"
	ErrReportNotFound                                  = "USER_REPORT_002"
	ErrReportUnauthorizedDelete                        = "USER_REPORT_003"
	ErrReportRequestPayloadValidationFailed            = "USER_REPORT_004"
	ErrReportRequestPayloadBadRequestJsonBadlyFormated = "USER_REPORT_005"

	// Infra errors
	ErrSaveTokenToDynamo = "DYNAMO_01"

	// Auth errors
	ErrAuthInvalidToken            = "AUTH_001"
	ErrAuthErrorDecodingHeader     = "AUTH_002"
	ErrAuthErrorDecodingPayload    = "AUTH_003"
	ErrAuthErrorUnmarshalingClaims = "AUTH_004"
	ErrAuthInvalidTokenSignature   = "AUTH_005"
)

// errorMessages defines a mapping of error codes to their corresponding descriptive error messages.
var errorMessages = map[string]string{
	ErrUserRequestPayloadBadRequestValidationError:   "Request payload is invalid, validation errors on inputs.",
	ErrUserRequestPayloadBadRequestEmptyFields:       "Request payload is invalid, required fields are empty.",
	ErrUserRequestPayloadBadRequestJsonBadlyFormated: "Request payload for user is invalid, JSON is badly formatted.",
	ErrRoleRequestPayloadBadRequestJsonBadlyFormated: "Request payload for role is invalid, JSON is badly formatted.",
	ErrRoleRequestPayloadBadRequestEmptyFields:       "Request payload for role is invalid, required fields are empty.",
	ErrRoleDatabaseUnableToCompleteOperation:         "Unable to complete operation on specified role.",
	ErrUserDatabaseUnableToCompleteOperation:         "Unable to complete operation on specified user.",
	ErrRoleRequestPayloadBadRequestValidationError:   "Request payload for role is invalid, validation errors on inputs.",
	ErrRoleNotFound:                                    "The specified role was not found: %v",
	ErrUserNotFound:                                    "The specified user was not found: %v",
	ErrRoleIncompleteParamsForCompleteUpdate:           "Incomplete role parameters for a full update.",
	ErrUserIncompleteParamsForCompleteUpdate:           "Incomplete user parameters for a full update.",
	ErrCreatingUserInCognitoUserPool:                   "Unable to create user in Cognito User Pool: %v",
	ErrSettingPermanentPassword:                        "Error setting permanent password: %v",
	ErrDeletingUserFromCognitoUserPool:                 "Error deleting user from Cognito: %v",
	ErrSaveTokenToDynamo:                               "Error saving token to DynamoDB: %v",
	ErrAuthInvalidToken:                                "Invalid JWT token.",
	ErrAuthInvalidTokenSignature:                       "Token signature is invalid.", // New error message
	ErrAuthErrorDecodingHeader:                         "Error decoding header: %v",
	ErrAuthErrorDecodingPayload:                        "Error decoding payload: %v",
	ErrAuthErrorUnmarshalingClaims:                     "Error unmarshaling claims: %v",
	ErrReportDatabaseUnableToCompleteOperation:         "Unable to complete operation on specified report.",
	ErrReportNotFound:                                  "Report not found with ID: %s",
	ErrReportUnauthorizedDelete:                        "The report does not belong to the provided user.",
	ErrReportRequestPayloadValidationFailed:            "Report request payload validation failed, missing or invalid fields.",
	ErrReportRequestPayloadBadRequestJsonBadlyFormated: "Report request payload is invalid, JSON is badly formatted.",
}

// GetErrorMessage returns a human-readable error message corresponding to the provided error code.
// If the error code does not match any predefined messages, it returns "unknown error".
func GetErrorMessage(code string) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "Unknown error."
}
