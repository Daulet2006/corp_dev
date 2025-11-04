package middleware

import (
	"cursed_backend/internal/logger"
	"net/http"
	"runtime/debug"

	"os"

	"cursed_backend/internal/models"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if e.Type == gin.ErrorTypePrivate {
					logger.Log.WithError(e.Err).Error("Internal error occurred")
					if os.Getenv("LOG_LEVEL") == "debug" {
						logger.Log.WithField("stack", string(debug.Stack())).Debug("Full stack trace")
					}
					c.JSON(http.StatusInternalServerError, models.APIResponse{
						Success: false,
						Message: "Internal server error",
					})
					return
				}
			}
		}
	}
}
