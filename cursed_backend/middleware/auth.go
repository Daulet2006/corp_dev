// Обновленный middleware/auth.go — добавим более гибкий RoleAuth для нескольких ролей
package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"cursed_backend/db"
	"cursed_backend/models"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func init() {
	secretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		secretKey = []byte("your-secret-key") // fallback для dev
	}
}

// JWTAuth middleware (без изменений, но добавим import os)
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)

		// Проверяем blocked
		var user models.User
		if err := db.GormDB.First(&user, userID).Error; err != nil || user.Blocked {
			c.JSON(http.StatusForbidden, gin.H{"error": "User blocked or not found"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}

// RoleAuth для одной роли (для admin-only)
func RoleAuth(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("%s role required", requiredRole)})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RoleOrAuth для нескольких ролей (e.g., manager or admin)
func RoleOr(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		allowed := false
		for _, r := range requiredRoles {
			if userRole == r {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role privileges"})
			c.Abort()
			return
		}
		c.Next()
	}
}
