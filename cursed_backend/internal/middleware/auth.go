package middleware

import (
	"cursed_backend/internal/db"
	"cursed_backend/internal/logger"
	"cursed_backend/internal/models"
	"fmt"
	"net/http"
	"strings"

	"cursed_backend/internal/security" // Импорт для ParseJWT

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := security.ParseJWT(tokenStr)
		if err != nil || claims == nil {
			logger.Log.WithError(err).Warn("Invalid token attempt")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)

		var user models.User
		if err := db.GormDB.First(&user, userID).Error; err != nil || user.Blocked {
			logger.AuditLog("auth_fail_blocked", userID, c.ClientIP(), err)
			c.JSON(http.StatusForbidden, gin.H{"error": "User blocked or not found"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}

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
