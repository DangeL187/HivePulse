package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"auth/internal/app"
)

func UserHasPermission(app *app.App, obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authUserIDValue, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		authUserID := authUserIDValue.(uint)
		authUserIDStr := strconv.FormatUint(uint64(authUserID), 10)

		allowed, err := app.RoleManager.Enforce(authUserIDStr, obj, act)
		if err != nil {
			zap.S().Errorf("Failed to enforce:\n%f", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check permission"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
