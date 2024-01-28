package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusForbidden, "missing token")
		}

		// Extract the token from the Authorization header
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusForbidden, "invalid token format")
		}
		tokenString = parts[1]

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Set the signing key here
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusForbidden, "invalid token")
		}

		// Check if it's a refresh token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && claims["type"] == "refresh" {
			// Handle refresh token logic here
			// Example: check if the refresh token is valid and generate a new access token
			// If refresh token is not valid or expired, return an error
			return echo.NewHTTPError(http.StatusForbidden, "refresh token is not allowed for this endpoint")
		}

		// Set user context for access token
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user", claims["user"])

		return next(c)
	}
}
