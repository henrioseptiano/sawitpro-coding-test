package middleware

import (
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestJWTMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "verysecret")
	jwtToken, _, _ := utils.GenerateJWTToken("mockUserID", "verysecret")
	tests := []struct {
		name             string
		token            string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Valid Token",
			token:            jwtToken,
			expectedStatus:   http.StatusOK,
			expectedResponse: "Authorized",
		},
		{
			name:             "Missing Token",
			token:            "",
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "missing token",
		},
		{
			name:             "Invalid Token Format",
			token:            "InvalidTokenFormat",
			expectedStatus:   http.StatusForbidden,
			expectedResponse: "invalid token format",
		},
		// Add more test cases for Expired Token, Refresh Token, etc.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup Echo context
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			req.Header.Set("Authorization", "Bearer "+tc.token)

			// Invoke middleware
			h := JWTMiddleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "Authorized")
			})
			err := h(c)
			if err != nil {
				assert.NotEqual(t, tc.expectedStatus, rec.Code)
			} else {
				// Assertions
				assert.Equal(t, tc.expectedStatus, rec.Code)
			}
		})
	}
}
