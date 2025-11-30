package middleware

import (
	"log"
	"net/http"

	"github.com/DangeL187/erax"
	"github.com/gin-gonic/gin"

	"auth/internal/app"
)

func Auth(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			err = erax.Wrap(err, "failed to find access_token cookie")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userID, tokenType, err := app.JWTManager.ParseToken(tokenString)
		if err != nil {
			err = erax.Wrap(err, "failed to parse token")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if tokenType != "access" {
			log.Printf("wrong token type (%s)\n", tokenType)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
