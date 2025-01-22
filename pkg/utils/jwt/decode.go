package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"styl-monolith/pkg/errorhandler"
)

type AccessJwtClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	TokenUse string `json:"token_use"`
	Expires  int64  `json:"exp"`
	IssuedAt int64  `json:"iat"`
}

type IdJwtClaims struct {
	Username  string `json:"cognito:username"`
	Email     string `json:"email"`
	TokenUse  string `json:"token_use"`
	Phone     string `json:"phone_number"`
	Expires   int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Birthdate string `json:"birthdate"`
}

// decodeJWT is a helper function that decodes a JWT token without validating its signature
// and unmarshals its payload into the provided claims object.
func decodeJWT(token string, claims interface{}) error {
	// Split the token into its parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errorhandler.NewDomainError(
			errorhandler.ErrAuthInvalidToken,
			errorhandler.GetErrorMessage(errorhandler.ErrAuthInvalidToken),
			nil)
	}

	// Decode and Validate Header (optional if not needed by business logic)
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

// verifyJWTSignature validates the given JWT token's signature using the key.
// It assumes HMAC-SHA256 signature for simplicity. You can extend this for other algorithms as required.
func verifyJWTSignature(token, key string) error {
	// Split token into its three parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid token format")
	}
	headerAndPayload := parts[0] + "." + parts[1]
	signature := parts[2]

	// Decode provided signature
	decodedSignature, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("failed to decode token signature")
	}

	// Generate HMAC-SHA256 signature using the header and payload with the signing key
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(headerAndPayload))
	expectedSignature := h.Sum(nil)

	// Compare generated signature with the provided signature
	if !hmac.Equal(decodedSignature, expectedSignature) {
		return errors.New("token signature validation failed")
	}

	return nil
}

// DecodeAndVerifyAccessJWT decodes an access JWT token into AccessJwtClaims and verifies its signature.
func DecodeAndVerifyAccessJWT(token, key string) (AccessJwtClaims, error) {
	var claims AccessJwtClaims

	// Verify the token's signature
	if err := verifyJWTSignature(token, key); err != nil {
		return claims, errorhandler.NewDomainError(
			errorhandler.ErrAuthInvalidTokenSignature,
			"Token signature validation failed",
			err)
	}

	// Decode the payload and map the claims
	err := decodeJWT(token, &claims)
	return claims, err
}

// DecodeAndVerifyIdJWT decodes an ID JWT token into IdJwtClaims and verifies its signature.
func DecodeAndVerifyIdJWT(token, key string) (IdJwtClaims, error) {
	var claims IdJwtClaims

	// Verify the token's signature
	if err := verifyJWTSignature(token, key); err != nil {
		return claims, errorhandler.NewDomainError(
			errorhandler.ErrAuthInvalidTokenSignature,
			"Token signature validation failed",
			err)
	}

	// Decode the payload and map the claims
	err := decodeJWT(token, &claims)
	return claims, err
}
