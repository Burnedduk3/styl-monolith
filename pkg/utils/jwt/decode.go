package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"strings"
	"styl-monolith/pkg/errorhandler"
)

type AccessJwtClaims struct {
	Username string              `json:"username"`
	Email    string              `json:"email"`
	TokenUse string              `json:"token_use"`
	Expires  timestamp.Timestamp `json:"exp"`
	IssuedAt timestamp.Timestamp `json:"iat"`
}
type IdJwtClaims struct {
	Username  string              `json:"cognito:username"`
	Email     string              `json:"email"`
	TokenUse  string              `json:"token_use"`
	Phone     string              `json:"phone_number"`
	Expires   timestamp.Timestamp `json:"exp"`
	IssuedAt  timestamp.Timestamp `json:"iat"`
	Birthdate string              `json:"birthdate"`
}

// decodeJWT is a helper function that decodes a JWT token and unmarshals its payload into the provided claims object.
func decodeJWT(token string, claims interface{}) error {
	parts := strings.Split(token, ".")
	if len(parts) < 3 {
		return errorhandler.NewDomainError(
			errorhandler.ErrAuthInvalidToken,
			errorhandler.GetErrorMessage(errorhandler.ErrAuthInvalidToken),
			nil)
	}

	// Decode Header
	if _, err := base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrAuthErrorDecodingHeader,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrAuthErrorDecodingHeader), err),
			err)
	}

	// Decode Payload
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrAuthErrorDecodingPayload,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrAuthErrorDecodingPayload), err),
			err)
	}

	// Unmarshal Payload JSON into Claims
	if err = json.Unmarshal(payload, claims); err != nil {
		return errorhandler.NewDomainError(
			errorhandler.ErrAuthErrorUnmarshalingClaims,
			fmt.Sprintf(errorhandler.GetErrorMessage(errorhandler.ErrAuthErrorUnmarshalingClaims), err),
			err)
	}

	return nil
}

// DecodeAccessJWT decodes an access JWT token into AccessJwtClaims.
func DecodeAccessJWT(token string) (AccessJwtClaims, error) {
	var claims AccessJwtClaims
	err := decodeJWT(token, &claims)
	return claims, err
}

// DecodeIdJWT decodes an ID JWT token into IdJwtClaims.
func DecodeIdJWT(token string) (IdJwtClaims, error) {
	var claims IdJwtClaims
	err := decodeJWT(token, &claims)
	return claims, err
}
