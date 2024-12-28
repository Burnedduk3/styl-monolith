package errorhandler

const (
	// User domain errorhandler
	ErrUserRequestPayloadBadRequest = "USER_001"

	// Order domain errorhandler
	ErrOrderNotFound     = "ORDER_001"
	ErrInvalidOrderState = "ORDER_002"
)

var errorMessages = map[string]string{
	ErrUserRequestPayloadBadRequest: "request payload is wrong",
	ErrOrderNotFound:                "order not found",
	ErrInvalidOrderState:            "invalid order state",
}

func GetErrorMessage(code string) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "unknown error"
}
