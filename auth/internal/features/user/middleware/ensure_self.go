package middleware

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"auth/internal/app"
	"auth/internal/shared/strutil"
)

func EnsureSelfByID(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		authUserIDValue, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		authUserID := authUserIDValue.(uint)

		userIDParam := c.Param("id")
		userID, err := strutil.StringToUint(userIDParam)
		if err != nil {
			zap.S().Errorf("Failed to parse user ID:\n%f", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		exists, err = app.UserModule.User.UserExists(authUserID)
		if err != nil {
			zap.S().Errorf("Failed to check if user exists:\n%f", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		if authUserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
