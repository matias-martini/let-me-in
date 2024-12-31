package auth 

import (
	"errors"
	"time"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"os"
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

func GenerateRefreshToken() (string, error) {
	// Generate a random 32-byte token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", errors.New("failed to generate refresh token")
	}

	// Encode to base64 for easy storage and transmission
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
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

// GenerateSalt creates a random salt string
func generateSalt() (string, error) {
	saltBytes := make([]byte, 16) // 16 bytes for salt
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(saltBytes), nil
}

// HashPassword hashes the password using salt and pepper
func hashPassword(password, salt string) (string, error) {
	pepper := os.Getenv("PEPPER")
	if pepper == "" {
		pepper = "PPR"
	}

	// Combine password, salt, and pepper
	combined := password + salt + pepper

	// Hash the combined password
	hashed, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword compares a plain text password with the hashed one
func verifyPassword(plainPassword, salt, hashedPassword string) bool {
	pepper := os.Getenv("PEPPER")
	if pepper == "" {
		pepper = "PPR"
	}

	// Combine password, salt, and pepper
	combined := plainPassword + salt + pepper

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(combined))
	return err == nil
}
