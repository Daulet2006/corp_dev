package security

import (
	"time"

	"cursed_backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// InitJWTSecret must be called after .env load (in main.go)
func InitJWTSecret(secret string) {
	if secret == "" {
		panic("JWT_SECRET is required")
	}
	jwtSecret = []byte(secret)
}

// GenerateJWT creates a JWT token for user
func GenerateJWT(user *models.User, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role.String(),
		"exp":     time.Now().Add(exp).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// GenerateJWTWithClaims for refresh/custom exp
func GenerateJWTWithClaims(userID uint, role string, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(exp).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// ParseJWT for validation (used in middleware)
func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
