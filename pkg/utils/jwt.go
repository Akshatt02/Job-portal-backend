package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a signed token using the provided secret.
// userID should be a string (UUID).
func GenerateJWT(userID, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken returns user_id string if token is valid.
func ParseToken(tokenStr, secret string) (string, error) {
	parser := &jwt.Parser{}
	token, err := parser.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Validate alg
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if uid, ok := claims["user_id"].(string); ok {
			return uid, nil
		}
	}

	return "", errors.New("invalid token claims")
}
