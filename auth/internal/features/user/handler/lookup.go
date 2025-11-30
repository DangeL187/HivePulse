package handler

import (
	"log"
	"net/http"

	"github.com/DangeL187/erax"
	"github.com/gin-gonic/gin"

	"auth/internal/app"
)

type UserEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func Lookup(app *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := UserEmailRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			err = erax.Wrap(err, "failed to parse user email request")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		userID, err := app.UserModule.User.GetUserIDByEmail(req.Email)
		if err != nil {
			err = erax.Wrap(err, "failed to get user ID by email")
			log.Printf("\n%f\n", err)

			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
		})
	}
}
