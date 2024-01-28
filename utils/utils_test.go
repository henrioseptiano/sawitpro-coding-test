package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid Phone Number", "+1234567890", true},
		{"Invalid Phone Number", "1234567890", false},
		{"Invalid Phone Number Format", "+1234567890a", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckPhoneNumber(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		expectedError bool
	}{
		{"Valid Password", "password123", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tc.password)
			assert.NoError(t, err)

			err = CheckPassword(tc.password, hashedPassword)
			assert.NoError(t, err)

			err = CheckPassword("wrongpassword", hashedPassword)
			assert.Error(t, err)

		})
	}
}

func TestGenerateAndDecodeJWTToken(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		secret      string
		expectedErr bool
	}{
		{"Valid Token Generation", "testUserID", "testSecret", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			accessToken, _, err := GenerateJWTToken(tc.userID, tc.secret)
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				userIDFromToken, err := DecodeJWTToken(accessToken)
				assert.NoError(t, err)
				assert.Equal(t, tc.userID, *userIDFromToken)
			}
		})
	}
}

func TestGetTokenFromAuthHeader(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		expected    string
		expectedErr bool
	}{
		{"Valid Authorization Header", "Bearer token123", "token123", false},
		{"Empty Authorization Header", "", "token123", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := GetTokenFromAuthHeader(tc.authHeader)
			if !tc.expectedErr {
				assert.Equal(t, tc.expected, token)
			} else {
				assert.NotEqual(t, tc.expected, token)
			}
		})
	}
}
