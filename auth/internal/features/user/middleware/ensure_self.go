package middleware

import (
	"log"
	"net/http"

	"github.com/DangeL187/erax"
	"github.com/gin-gonic/gin"

	"auth/internal/app"
	"auth/internal/shared/strutil"
)

func EnsureSelfByEmail(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		authUserIDValue, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		authUserID := authUserIDValue.(uint)

		userEmail := c.Param("email")
		if userEmail == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}

		userID, err := app.UserModule.User.GetUserIDByEmail(userEmail)
		if err != nil {
			err = erax.Wrap(err, "failed to get user ID by email")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		if authUserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

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
			err = erax.Wrap(err, "failed to parse user ID")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		exists, err = app.UserModule.User.UserExists(authUserID)
		if err != nil {
			err = erax.Wrap(err, "failed to check if user exists")
			log.Printf("\n%f\n", err)

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
