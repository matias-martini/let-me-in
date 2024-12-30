package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key (use environment variables in production)
var jwtSecret = []byte("super-secret-key")

// Claims structure
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token for a user
func GenerateJWT(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), // Token expires in 2 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates and parses a JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

