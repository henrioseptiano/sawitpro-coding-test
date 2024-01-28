package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func CheckPhoneNumber(s string) bool {
	if strings.HasPrefix(s, "+") {
		_, err := strconv.Atoi(s[1:])
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func HashPassword(password string) (string, error) {
	cost := 10

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(providedPassword, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func GenerateJWTToken(userID, secret string) (string, string, error) {
	signingKey := []byte(secret)
	// Generate access token
	accessTokenClaims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // Access token expires in 15 minutes
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(signingKey)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshTokenClaims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // Refresh token expires in 7 days
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(signingKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func DecodeJWTToken(tokenString string) (*string, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	// Access the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	// Extract specific claims
	username, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("User Id Not Found")
	}

	return &username, nil
}

func GetTokenFromAuthHeader(auth string) string {
	strs := strings.Split(auth, " ")
	if len(strs) > 1 {
		return strs[1]
	}
	return ""
}
