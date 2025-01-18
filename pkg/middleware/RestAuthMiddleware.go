package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"styl-monolith/internal/users/core/ports"
	"styl-monolith/pkg/utils/jwt"
)

func ValidateTokenWithRepository(repo ports.LoginPort) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token hash from the request header
			tokenHash := c.Request().Header.Get("Authorization")
			if tokenHash == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing token hash in Authorization header")
			}
			// Read and parse the body
			bodyBytes, err := io.ReadAll(c.Request().Body) // Reading the request body
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Failed to read request body")
			}

			// Restore the body stream so other handlers/middleware can read it
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Parse the body as JSON to extract the email (if expected format is JSON)
			var bodyData map[string]interface{}
			if err = json.Unmarshal(bodyBytes, &bodyData); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON format in request body")
			}

			email, ok := bodyData["email"].(string)
			if !ok || email == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Email not found in request body")
			}

			// Query the repository to get token data by token ID
			token, err := repo.GetTokensFromDynamoById(tokenHash, email)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Failed to retrieve token: "+err.Error())
			}

			// Decode and validate the token
			claims, err := jwt.DecodeAccessJWT(token.AccessToken) // Assuming DecodeAccessJWT function exists
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			// Attach claims to context for downstream usage
			c.Set("userClaims", claims)

			// Pass control to the next handler
			return next(c)
		}
	}
}
