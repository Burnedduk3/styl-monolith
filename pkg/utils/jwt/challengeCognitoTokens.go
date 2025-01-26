package jwt

import "errors"

func ValidateCognitoTokens(accessToken, idToken string) error {
	// Decode and validate Access Token
	accessTokenClaims, err := DecodeAndVerifyAccessJWT(accessToken, idToken)
	if err != nil {
		return errors.New("user not authenticated: invalid access token")
	}

	// Decode and validate ID Token
	idTokenClaims, err := DecodeAndVerifyIdJWT(idToken, idToken)
	if err != nil || idTokenClaims.Email != accessTokenClaims.Email {
		return errors.New("user not authenticated: invalid or mismatched id token")
	}

	// Map user data and claims for return
	return nil
}
