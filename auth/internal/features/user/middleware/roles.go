package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DangeL187/erax"
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
			err = erax.Wrap(err, "failed to enforce")
			log.Printf("\n%f\n", err)

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
