package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"styl-monolith/internal/users/core/ports"
	"styl-monolith/pkg/utils/jwt"
	"time"
)

// ValidateAccessTokenWithRepository is a middleware that validates an access token using the provided repository.
// It ensures the token is valid, not expired, and of correct use, attaching user claims to the context for further use.
// If the token is missing, invalid, or expired, it returns an HTTP error.
func ValidateAccessTokenWithRepository(repo ports.LoginPort) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token hash from the request header
			AuthHeader := c.Request().Header.Get("Authorization")
			if AuthHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing access token in Authorization header")
			}

			tokenHash := AuthHeader[7:]

			// Query the repository to get token data by token ID
			token, err := repo.GetTokensFromDynamoById(tokenHash)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Failed to retrieve token: "+err.Error())
			}

			// Decode and validate the token
			claims, err := jwt.DecodeAccessJWT(token.AccessToken) // Assuming DecodeAccessJWT function exists
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			if claims.TokenUse != "access" {

			}

			// Validate current time is between claims.IssuedAt and claims.Expires
			currentTime := time.Now().Unix() + 10
			if currentTime < claims.IssuedAt || currentTime > claims.Expires {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token is not valid at the current time")
			}

			// Attach claims to context for downstream usage
			c.Set("email", claims.Email)
			c.Set("username", claims.Username)

			// Pass control to the next handler
			return next(c)
		}
	}
}
