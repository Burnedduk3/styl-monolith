package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"styl-monolith/internal/users/core/ports"
	"styl-monolith/pkg/utils/jwt"
	"styl-monolith/pkg/validator"
)

// ValidateAccessTokenWithRepository is a middleware that validates an access token using the provided repository.
// It ensures the token is valid, not expired, and of correct use, attaching user claims to the context for further use.
// If the token is missing, invalid, or expired, it returns an HTTP error.
func ValidateAccessTokenWithRepository(repo ports.LoginPort) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the Authorization header
			AuthHeader := c.Request().Header.Get("Authorization")
			if AuthHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			// Ensure the token starts with Bearer
			if !strings.HasPrefix(AuthHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Malformed Authorization header")
			}

			// Extract the token after "Bearer "
			tokenString := AuthHeader[7:]
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing token after Bearer")
			}

			// Query the repository to ensure the token exists
			token, err := repo.GetTokensFromDynamoById(tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Failed to retrieve token: "+err.Error())
			}

			// Verify the token's signature and decode it
			claims, err := jwt.DecodeAndVerifyAccessJWT(token.AccessToken, token.IdTokenHash)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token signature or expired token")
			}

			// Check the token use (ensure it's specifically for 'access')
			if claims.TokenUse != "access" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token use")
			}

			// Validate the token timestamps
			if !validator.ValidateIfActualTimeIsBetweenTwoTimestamps(claims.IssuedAt, claims.Expires) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired or is not yet valid")
			}

			// Optional: Ensure token has required claims (e.g., Email, Username)
			if claims.Email == "" || claims.Username == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token is missing required claims")
			}

			// Attach validated claims to the context for subsequent use
			c.Set("email", claims.Email)
			c.Set("username", claims.Username)

			// Pass control to the next handler in the chain
			return next(c)

		}
	}
}
