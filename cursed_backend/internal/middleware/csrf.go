package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"cursed_backend/internal/logger"

	"github.com/gin-gonic/gin"
)

const csrfCookieName = "csrf_token"

// CSRF middleware: Manual double-submit cookie for protected mutating requests
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get token from header
		clientToken := c.GetHeader("X-CSRF-Token")
		if clientToken == "" {
			writeCSRFError(c, "CSRF token missing")
			return
		}

		// Get cookie token
		cookie, err := c.Request.Cookie(csrfCookieName)
		if err != nil || cookie.Value == "" {
			writeCSRFError(c, "CSRF cookie missing")
			return
		}

		// Validate match
		if clientToken != cookie.Value {
			userID := uint(0)
			if uid, _ := c.Get("user_id"); uid != nil {
				userID = uint(uid.(float64))
			}
			logger.Log.WithFields(map[string]interface{}{
				"path":        c.Request.URL.Path,
				"ip":          c.ClientIP(),
				"user_id":     userID,
				"error":       "CSRF token mismatch",
				"clientToken": clientToken[:8] + "...", // Partial log
				"cookieToken": cookie.Value[:8] + "...",
			}).Warn("CSRF validation failed")
			writeCSRFError(c, "CSRF token invalid")
			return
		}

		// Optional: Regenerate token after use (for extra security, but complicates)
		// newToken, _ := generateToken()
		// setCSRF(c, newToken)
		// c.Header("X-CSRF-Token-New", newToken) // Client rotate

		c.Next()
	}
}

// Generate random 32-byte token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Set httpOnly secure cookie
func setCSRF(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   3600 * 24, // 24h
		HttpOnly: true,
		Secure:   c.Request.TLS != nil, // Prod HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

// Write 403 error + abort
func writeCSRFError(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, gin.H{"error": msg})
	c.Abort()
}

// CSRFToken handler: Generate & set cookie, return token
func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := generateToken()
		if err != nil {
			logger.Log.WithError(err).Error("Failed to generate CSRF token")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		// Set cookie
		setCSRF(c, token)

		// Return token
		c.JSON(http.StatusOK, gin.H{"csrf_token": token})

		userID := uint(0)
		if uid, _ := c.Get("user_id"); uid != nil {
			userID = uint(uid.(float64))
		}
		logger.Log.WithField("user_id", userID).Info("CSRF token generated and cookie set")
	}
}
